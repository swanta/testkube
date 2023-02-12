package slack

import (
	"bytes"
	"encoding/json"
	"os"
	"text/template"

	"github.com/slack-go/slack"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/log"
)

type MessageArgs struct {
	ExecutionName string
	EventType     string
	Namespace     string
	Labels        string
	TestName      string
	TestType      string
	Status        string
	FailedSteps   int
	TotalSteps    int
	StartTime     string
	EndTime       string
	Duration      string
}

type Notifier struct {
	client          *slack.Client
	timestamps      map[string]string
	Ready           bool
	messageTemplate string
	config          *Config
}

func NewNotifier(template string, config []NotificationsConfig) *Notifier {
	notifier := Notifier{messageTemplate: template, config: NewConfig(config)}
	notifier.timestamps = make(map[string]string)
	if token, ok := os.LookupEnv("SLACK_TOKEN"); ok {
		log.DefaultLogger.Infow("initializing slack client", "SLACK_TOKEN", token)
		notifier.client = slack.New(token, slack.OptionDebug(true))
		notifier.Ready = true
	} else {
		log.DefaultLogger.Warn("SLACK_TOKEN is not set")
	}
	return &notifier
}

// SendMessage posts a message to the slack configured channel
func (s *Notifier) SendMessage(channelID string, message string) error {
	if s.client != nil {
		_, _, err := s.client.PostMessage(channelID, slack.MsgOptionText(message, false))
		if err != nil {
			log.DefaultLogger.Warnw("error while posting message to channel", "channelID", channelID, "error", err.Error())
			return err
		}
	} else {
		log.DefaultLogger.Warnw("slack client is not initialised")
	}
	return nil
}

// SendEvent composes an event message and sends it to slack
func (s *Notifier) SendEvent(event *testkube.Event) error {

	message, name, err := s.composeMessage(event)
	if err != nil {
		return err
	}

	if s.client != nil {

		log.DefaultLogger.Debugw("sending event to slack", "event", event)
		channels, err := s.getChannels(event)
		if err != nil {
			return err
		}
		log.DefaultLogger.Infow("channels to send event to", "channels", channels)

		for _, channelID := range channels {
			prevTimestamp, ok := s.timestamps[name]
			var timestamp string

			if ok {
				_, timestamp, _, err = s.client.UpdateMessage(channelID, prevTimestamp, slack.MsgOptionBlocks(message.Blocks.BlockSet...))
			}

			if !ok || err != nil {
				_, timestamp, err = s.client.PostMessage(channelID, slack.MsgOptionBlocks(message.Blocks.BlockSet...))
			}

			if err != nil {
				log.DefaultLogger.Warnw("error while posting message to channel", "channelID", channelID, "error", err.Error())
				return err
			}

			if event.IsSuccess() {
				delete(s.timestamps, name)
			} else {
				s.timestamps[name] = timestamp
			}
		}
	} else {
		log.DefaultLogger.Warnw("slack client is not initialised")
	}

	return nil
}

func (s *Notifier) getChannels(event *testkube.Event) ([]string, error) {
	result := []string{}
	if !s.config.HasChannelsDefined() {
		channels, _, err := s.client.GetConversationsForUser(&slack.GetConversationsForUserParameters{})
		if err != nil {
			log.DefaultLogger.Warnw("error while getting bot channels", "error", err.Error())
			return nil, err
		}
		_, needsSending := s.config.NeedsSending(event)
		if len(channels) > 0 && needsSending {
			result = append(result, channels[0].GroupConversation.ID)
			return result, nil
		}
	} else {
		channels, needsSending := s.config.NeedsSending(event)
		if needsSending {
			return channels, nil
		}
	}
	return nil, nil
}

func (s *Notifier) composeMessage(event *testkube.Event) (view *slack.Message, name string, err error) {
	var message []byte
	if event.TestExecution != nil {
		message, err = s.composeTestMessage(event.TestExecution, event.Type())
		name = event.TestExecution.Name
	} else if event.TestSuiteExecution != nil {
		message, err = s.composeTestsuiteMessage(event.TestSuiteExecution, event.Type())
		name = event.TestSuiteExecution.Name
	} else {
		log.DefaultLogger.Warnw("event type is not handled by Slack notifier", "event", event)
		return nil, "", nil
	}

	if err != nil {
		return nil, "", err
	}
	view = &slack.Message{}
	err = json.Unmarshal(message, view)
	if err != nil {
		log.DefaultLogger.Warnw("error while creating slack specific message", "error", err.Error())
		return nil, "", err
	}

	return view, name, nil
}

func (s *Notifier) composeTestsuiteMessage(execution *testkube.TestSuiteExecution, eventType testkube.EventType) ([]byte, error) {
	t, err := template.New("message").Parse(s.messageTemplate)
	if err != nil {
		log.DefaultLogger.Warnw("error while parsing slack template", "error", err.Error())
		return nil, err
	}

	args := MessageArgs{
		ExecutionName: execution.Name,
		EventType:     string(eventType),
		Namespace:     execution.TestSuite.Namespace,
		Labels:        testkube.MapToString(execution.Labels),
		TestName:      execution.TestSuite.Name,
		TestType:      "Test Suite",
		Status:        string(*execution.Status),
		StartTime:     execution.StartTime.String(),
		EndTime:       execution.EndTime.String(),
		Duration:      execution.Duration,
		TotalSteps:    len(execution.StepResults),
		FailedSteps:   execution.FailedStepsCount(),
	}

	log.DefaultLogger.Infow("Execution changed", "status", execution.Status)

	var message bytes.Buffer
	err = t.Execute(&message, args)
	if err != nil {
		log.DefaultLogger.Warnw("error while executing slack template", "error", err.Error())
		return nil, err
	}
	return message.Bytes(), nil
}

func (s *Notifier) composeTestMessage(execution *testkube.Execution, eventType testkube.EventType) ([]byte, error) {
	t, err := template.New("message").Parse(s.messageTemplate)
	if err != nil {
		log.DefaultLogger.Warnw("error while parsing slack template", "error", err.Error())
		return nil, err
	}

	args := MessageArgs{
		ExecutionName: execution.Name,
		EventType:     string(eventType),
		Namespace:     execution.TestNamespace,
		Labels:        testkube.MapToString(execution.Labels),
		TestName:      execution.TestName,
		TestType:      execution.TestType,
		Status:        string(*execution.ExecutionResult.Status),
		StartTime:     execution.StartTime.String(),
		EndTime:       execution.EndTime.String(),
		Duration:      execution.Duration,
		TotalSteps:    len(execution.ExecutionResult.Steps),
		FailedSteps:   execution.ExecutionResult.FailedStepsCount(),
	}

	log.DefaultLogger.Infow("Execution changed", "status", execution.ExecutionResult.Status)

	var message bytes.Buffer
	err = t.Execute(&message, args)
	if err != nil {
		log.DefaultLogger.Warnw("error while executing slack template", "error", err.Error())
		return nil, err
	}
	return message.Bytes(), nil
}

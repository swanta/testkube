package testresult

import (
	"context"
	"strings"
	"time"

	"github.com/kubeshop/testkube/pkg/repository/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
)

var _ Repository = &MongoRepository{}

const CollectionName = "testresults"

func NewMongoRepository(db *mongo.Database, allowDiskUse bool, opts ...MongoRepositoryOpt) *MongoRepository {
	r := &MongoRepository{
		Coll:         db.Collection(CollectionName),
		allowDiskUse: allowDiskUse,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

type MongoRepository struct {
	Coll         *mongo.Collection
	allowDiskUse bool
}

func WithMongoRepositoryCollection(collection *mongo.Collection) MongoRepositoryOpt {
	return func(r *MongoRepository) {
		r.Coll = collection
	}
}

type MongoRepositoryOpt func(*MongoRepository)

func (r *MongoRepository) Get(ctx context.Context, id string) (result testkube.TestSuiteExecution, err error) {
	err = r.Coll.FindOne(ctx, bson.M{"$or": bson.A{bson.M{"id": id}, bson.M{"name": id}}}).Decode(&result)
	return
}

func (r *MongoRepository) GetByNameAndTestSuite(ctx context.Context, name, testSuiteName string) (result testkube.TestSuiteExecution, err error) {
	err = r.Coll.FindOne(ctx, bson.M{"name": name, "testsuite.name": testSuiteName}).Decode(&result)
	return
}

func (r *MongoRepository) GetLatestByTestSuite(ctx context.Context, testSuiteName, sortField string) (result testkube.TestSuiteExecution, err error) {
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: sortField, Value: -1}})
	err = r.Coll.FindOne(ctx, bson.M{"testsuite.name": testSuiteName}, findOptions).Decode(&result)
	return
}

func (r *MongoRepository) GetLatestByTestSuites(ctx context.Context, testSuiteNames []string, sortField string) (executions []testkube.TestSuiteExecution, err error) {
	var results []struct {
		LatestID string `bson:"latest_id"`
	}

	if len(testSuiteNames) == 0 {
		return executions, nil
	}

	conditions := bson.A{}
	for _, testSuiteName := range testSuiteNames {
		conditions = append(conditions, bson.M{"testsuite.name": testSuiteName})
	}

	pipeline := []bson.D{{{Key: "$match", Value: bson.M{"$or": conditions}}}}
	pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: sortField, Value: -1}}}})
	pipeline = append(pipeline, bson.D{
		{Key: "$group", Value: bson.D{{Key: "_id", Value: "$testsuite.name"}, {Key: "latest_id", Value: bson.D{{Key: "$first", Value: "$id"}}}}}})

	optsA := options.Aggregate()
	if r.allowDiskUse {
		optsA.SetAllowDiskUse(r.allowDiskUse)
	}

	cursor, err := r.Coll.Aggregate(ctx, pipeline, optsA)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return executions, nil
	}

	conditions = bson.A{}
	for _, result := range results {
		conditions = append(conditions, bson.M{"id": result.LatestID})
	}

	optsF := options.Find()
	if r.allowDiskUse {
		optsF.SetAllowDiskUse(r.allowDiskUse)
	}

	cursor, err = r.Coll.Find(ctx, bson.M{"$or": conditions}, optsF)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &executions)
	if err != nil {
		return nil, err
	}

	return executions, nil
}

func (r *MongoRepository) GetNewestExecutions(ctx context.Context, limit int) (result []testkube.TestSuiteExecution, err error) {
	result = make([]testkube.TestSuiteExecution, 0)
	resultLimit := int64(limit)
	opts := &options.FindOptions{Limit: &resultLimit}
	opts.SetSort(bson.D{{Key: "_id", Value: -1}})
	if r.allowDiskUse {
		opts.SetAllowDiskUse(r.allowDiskUse)
	}

	cursor, err := r.Coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return result, err
	}
	err = cursor.All(ctx, &result)
	return
}

func (r *MongoRepository) GetExecutionsTotals(ctx context.Context, filter ...Filter) (totals testkube.ExecutionsTotals, err error) {
	var result []struct {
		Status string `bson:"_id"`
		Count  int    `bson:"count"`
	}

	query := bson.M{}
	if len(filter) > 0 {
		query, _ = composeQueryAndOpts(filter[0])
	}

	pipeline := []bson.D{{{Key: "$match", Value: query}}}
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "starttime", Value: -1}}}})
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: int64(filter[0].Page() * filter[0].PageSize())}})
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: int64(filter[0].PageSize())}})
	}

	pipeline = append(pipeline, bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$status"},
		{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}})

	opts := options.Aggregate()
	if r.allowDiskUse {
		opts.SetAllowDiskUse(r.allowDiskUse)
	}

	cursor, err := r.Coll.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return totals, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		return totals, err
	}

	var sum int32

	for _, o := range result {
		sum += int32(o.Count)
		switch testkube.TestSuiteExecutionStatus(o.Status) {
		case testkube.QUEUED_TestSuiteExecutionStatus:
			totals.Queued = int32(o.Count)
		case testkube.RUNNING_TestSuiteExecutionStatus:
			totals.Running = int32(o.Count)
		case testkube.PASSED_TestSuiteExecutionStatus:
			totals.Passed = int32(o.Count)
		case testkube.FAILED_TestSuiteExecutionStatus:
			totals.Failed = int32(o.Count)
		}
	}
	totals.Results = sum

	return
}

func (r *MongoRepository) GetExecutions(ctx context.Context, filter Filter) (result []testkube.TestSuiteExecution, err error) {
	result = make([]testkube.TestSuiteExecution, 0)
	query, opts := composeQueryAndOpts(filter)
	if r.allowDiskUse {
		opts.SetAllowDiskUse(r.allowDiskUse)
	}

	cursor, err := r.Coll.Find(ctx, query, opts)
	if err != nil {
		return
	}
	err = cursor.All(ctx, &result)
	return
}

func (r *MongoRepository) Insert(ctx context.Context, result testkube.TestSuiteExecution) (err error) {
	_, err = r.Coll.InsertOne(ctx, result)
	return
}

func (r *MongoRepository) Update(ctx context.Context, result testkube.TestSuiteExecution) (err error) {
	_, err = r.Coll.ReplaceOne(ctx, bson.M{"id": result.Id}, result)
	return
}

// StartExecution updates execution start time
func (r *MongoRepository) StartExecution(ctx context.Context, id string, startTime time.Time) (err error) {
	_, err = r.Coll.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"starttime": startTime}})
	return
}

// EndExecution updates execution end time
func (r *MongoRepository) EndExecution(ctx context.Context, e testkube.TestSuiteExecution) (err error) {
	_, err = r.Coll.UpdateOne(ctx, bson.M{"id": e.Id}, bson.M{"$set": bson.M{"endtime": e.EndTime, "duration": e.Duration, "durationms": e.DurationMs}})
	return
}

func composeQueryAndOpts(filter Filter) (bson.M, *options.FindOptions) {

	query := bson.M{}
	opts := options.Find()
	startTimeQuery := bson.M{}

	if filter.NameDefined() {
		query["testsuite.name"] = filter.Name()
	}

	if filter.TextSearchDefined() {
		query["name"] = bson.M{"$regex": primitive.Regex{Pattern: filter.TextSearch(), Options: "i"}}
	}

	if filter.LastNDaysDefined() {
		startTimeQuery["$gte"] = time.Now().Add(-time.Duration(filter.LastNDays()) * 24 * time.Hour)
	}

	if filter.StartDateDefined() {
		startTimeQuery["$gte"] = filter.StartDate()
	}

	if filter.EndDateDefined() {
		startTimeQuery["$lte"] = filter.EndDate()
	}

	if len(startTimeQuery) > 0 {
		query["starttime"] = startTimeQuery
	}

	if filter.StatusesDefined() {
		statuses := filter.Statuses()
		if len(statuses) == 1 {
			query["status"] = statuses[0]
		} else {
			var conditions bson.A
			for _, status := range statuses {
				conditions = append(conditions, bson.M{"status": status})
			}

			query["$or"] = conditions
		}
	}

	if filter.Selector() != "" {
		items := strings.Split(filter.Selector(), ",")
		for _, item := range items {
			elements := strings.Split(item, "=")
			if len(elements) == 2 {
				query["labels."+elements[0]] = elements[1]
			} else if len(elements) == 1 {
				query["labels."+elements[0]] = bson.M{"$exists": true}
			}
		}
	}

	opts.SetSkip(int64(filter.Page() * filter.PageSize()))
	opts.SetLimit(int64(filter.PageSize()))
	opts.SetSort(bson.D{{Key: "starttime", Value: -1}})

	return query, opts
}

// DeleteByTestSuite deletes execution results by test suite
func (r *MongoRepository) DeleteByTestSuite(ctx context.Context, testSuiteName string) (err error) {
	_, err = r.Coll.DeleteMany(ctx, bson.M{"testsuite.name": testSuiteName})
	return
}

// DeleteAll deletes all execution results
func (r *MongoRepository) DeleteAll(ctx context.Context) (err error) {
	_, err = r.Coll.DeleteMany(ctx, bson.M{})
	return
}

// DeleteByTestSuites deletes execution results by test suites
func (r *MongoRepository) DeleteByTestSuites(ctx context.Context, testSuiteNames []string) (err error) {
	if len(testSuiteNames) == 0 {
		return nil
	}

	var filter bson.M
	if len(testSuiteNames) > 1 {
		conditions := bson.A{}
		for _, testSuiteName := range testSuiteNames {
			conditions = append(conditions, bson.M{"testsuite.name": testSuiteName})
		}

		filter = bson.M{"$or": conditions}
	} else {
		filter = bson.M{"testsuite.name": testSuiteNames[0]}
	}

	_, err = r.Coll.DeleteMany(ctx, filter)
	return
}

// GetTestSuiteMetrics returns test executions metrics
func (r *MongoRepository) GetTestSuiteMetrics(ctx context.Context, name string, limit, last int) (metrics testkube.ExecutionsMetrics, err error) {
	query := bson.M{"testsuite.name": name}

	var pipeline []bson.D
	if last > 0 {
		query["starttime"] = bson.M{"$gte": time.Now().Add(-time.Duration(last) * 24 * time.Hour)}
	}

	pipeline = append(pipeline, bson.D{{Key: "$match", Value: query}})
	pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "starttime", Value: -1}}}})
	pipeline = append(pipeline, bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "status", Value: 1},
				{Key: "duration", Value: 1},
				{Key: "starttime", Value: 1},
				{Key: "name", Value: 1},
			},
		},
	})

	opts := options.Aggregate()
	if r.allowDiskUse {
		opts.SetAllowDiskUse(r.allowDiskUse)
	}

	cursor, err := r.Coll.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return metrics, err
	}

	var executions []testkube.ExecutionsMetricsExecutions
	err = cursor.All(ctx, &executions)

	if err != nil {
		return metrics, err
	}

	metrics = common.CalculateMetrics(executions)
	if limit > 0 && limit < len(metrics.Executions) {
		metrics.Executions = metrics.Executions[:limit]
	}

	return metrics, nil
}

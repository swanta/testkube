"use strict";(self.webpackChunktestkube_documentation=self.webpackChunktestkube_documentation||[]).push([[1254],{3905:(e,t,r)=>{r.d(t,{Zo:()=>u,kt:()=>d});var n=r(67294);function a(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function i(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function o(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?i(Object(r),!0).forEach((function(t){a(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):i(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function s(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},i=Object.keys(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var l=n.createContext({}),c=function(e){var t=n.useContext(l),r=t;return e&&(r="function"==typeof e?e(t):o(o({},t),e)),r},u=function(e){var t=c(e.components);return n.createElement(l.Provider,{value:t},e.children)},p={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},f=n.forwardRef((function(e,t){var r=e.components,a=e.mdxType,i=e.originalType,l=e.parentName,u=s(e,["components","mdxType","originalType","parentName"]),f=c(r),d=a,m=f["".concat(l,".").concat(d)]||f[d]||p[d]||i;return r?n.createElement(m,o(o({ref:t},u),{},{components:r})):n.createElement(m,o({ref:t},u))}));function d(e,t){var r=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var i=r.length,o=new Array(i);o[0]=f;var s={};for(var l in t)hasOwnProperty.call(t,l)&&(s[l]=t[l]);s.originalType=e,s.mdxType="string"==typeof e?e:a,o[1]=s;for(var c=2;c<i;c++)o[c]=r[c];return n.createElement.apply(null,o)}return n.createElement.apply(null,r)}f.displayName="MDXCreateElement"},18478:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>l,contentTitle:()=>o,default:()=>p,frontMatter:()=>i,metadata:()=>s,toc:()=>c});var n=r(87462),a=(r(67294),r(3905));const i={},o=void 0,s={unversionedId:"cli-reference/testkube_run_test",id:"cli-reference/testkube_run_test",title:"testkube_run_test",description:"testkube run test",source:"@site/docs/5-cli-reference/testkube_run_test.md",sourceDirName:"5-cli-reference",slug:"/cli-reference/testkube_run_test",permalink:"/testkube/cli-reference/testkube_run_test",draft:!1,editUrl:"https://github.com/kubeshop/testkube/docs/docs/5-cli-reference/testkube_run_test.md",tags:[],version:"current",frontMatter:{},sidebar:"tutorialSidebar",previous:{title:"testkube_run",permalink:"/testkube/cli-reference/testkube_run"},next:{title:"testkube_run_testsuite",permalink:"/testkube/cli-reference/testkube_run_testsuite"}},l={},c=[{value:"testkube run test",id:"testkube-run-test",level:2},{value:"Synopsis",id:"synopsis",level:3},{value:"Options",id:"options",level:3},{value:"Options inherited from parent commands",id:"options-inherited-from-parent-commands",level:3},{value:"SEE ALSO",id:"see-also",level:3}],u={toc:c};function p(e){let{components:t,...r}=e;return(0,a.kt)("wrapper",(0,n.Z)({},u,r,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h2",{id:"testkube-run-test"},"testkube run test"),(0,a.kt)("p",null,"Starts new test"),(0,a.kt)("h3",{id:"synopsis"},"Synopsis"),(0,a.kt)("p",null,"Starts new test based on Test Custom Resource name, returns results to console"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"testkube run test <testName> [flags]\n")),(0,a.kt)("h3",{id:"options"},"Options"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},'      --args stringArray                           executor binary additional arguments\n      --artifact-dir stringArray                   artifact dirs for container executor\n      --artifact-storage-class-name string         artifact storage class name for container executor\n      --artifact-volume-mount-path string          artifact volume mount path for container executor\n      --concurrency int                            concurrency level for multiple test execution (default 10)\n      --copy-files stringArray                     file path mappings from host to pod of form source:destination\n  -d, --download-artifacts                         downlaod artifacts automatically\n      --download-dir string                        download dir (default "artifacts")\n      --env stringToString                         envs in a form of name1=val1 passed to executor (default [])\n      --execution-label stringToString             execution-label key value pair: --execution-label key1=value1 (default [])\n      --git-branch string                          if uri is git repository we can set additional branch parameter\n      --git-commit string                          if uri is git repository we can use commit id (sha) parameter\n      --git-path string                            if repository is big we need to define additional path to directory/file to checkout partially\n      --git-working-dir string                     if repository contains multiple directories with tests (like monorepo) and one starting directory we can set working directory parameter\n  -h, --help                                       help for test\n      --http-proxy string                          http proxy for executor containers\n      --https-proxy string                         https proxy for executor containers\n      --image string                               execution variable passed to executor\n      --iterations int                             how many times to run the test (default 1)\n      --job-template string                        job template file path for extensions to job template\n  -l, --label strings                              label key value pair: --label key1=value1\n  -n, --name string                                execution name, if empty will be autogenerated\n      --negative-test                              negative test, if enabled, makes failure an expected and correct test result. If the test fails the result will be set to success, and vice versa\n      --prerun-script string                       path to script to be run before test execution\n      --scraper-template string                    scraper template file path for extensions to scraper template\n      --secret stringToString                      secret envs in a form of secret_key1=secret_name1 passed to executor (default [])\n  -s, --secret-variable stringToString             execution secret variable passed to executor (default [])\n      --secret-variable-reference stringToString   secret variable references in a form name1=secret_name1=secret_key1 (default [])\n  -v, --variable stringToString                    execution variable passed to executor (default [])\n      --variables-file string                      variables file path, e.g. postman env file - will be passed to executor if supported\n  -f, --watch                                      watch for changes after start\n')),(0,a.kt)("h3",{id:"options-inherited-from-parent-commands"},"Options inherited from parent commands"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},'  -a, --api-uri string     api uri, default value read from config if set (default "http://localhost:8088")\n  -c, --client string      client used for connecting to Testkube API one of proxy|direct (default "proxy")\n      --namespace string   Kubernetes namespace, default value read from config if set (default "testkube")\n      --oauth-enabled      enable oauth (default true)\n      --verbose            show additional debug messages\n')),(0,a.kt)("h3",{id:"see-also"},"SEE ALSO"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"/testkube/cli-reference/testkube_run"},"testkube run"),"\t - Runs tests or test suites")))}p.isMDXComponent=!0}}]);
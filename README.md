# vote microservice
Frontend service that renders an HTML page to vote on cats vs. dogs.
![cat-v-dog](https://user-images.githubusercontent.com/879348/95268443-3357cb00-07ec-11eb-8913-d83e322d26f0.png)

The "vote" service forwards requests to the ["api"](https://github.com/copilot-example-voting-app/api) microservice to store
and retrieve results on whether a voter prefers cats or dogs.

The two services communicate through Service Discovery which AWS Copilot sets up by default by querying the `api.voting-app.local:8080` endpoint.  
Alternatively, you can use the [`COPILOT_SERVICE_DISCOVERY_ENDPOINT` environment variable](https://github.com/copilot-example-voting-app/vote/blob/6b4a2dab38229b89e84d1aca6081a0577e9be167/server/server.go#L122) 
that Copilot injects by default to your service.

## How to create this service?
1. Install the AWS Copilot CLI [https://aws.github.io/copilot-cli/](https://aws.github.io/copilot-cli/)
2. Run
   ```bash
   $ copilot init
   ```
3. Enter "voting-app" for the name of your application.
4. Select "Load Balanced Web Service" for the service type.
5. Enter "vote" for the name of the service.
6. Say "Y" to deploying to a "test" environment 🚀

Once deployed, your service will be accessible at an HTTP endpoint provided by the CLI like: http://votin-publi-anelun2kxbrl-XXXXXXX.YYYYY.elb.amazonaws.com/

## What does it do?
AWS Copilot uses AWS CloudFormation under the hood to provision your infrastructure resources.
You should be able to see a `voting-app-test-vote` stack that yours ECS service along with all the peripheral resources
needed for logging, service discovery, and more...

## How does it work?
Copilot stores the infrastructure-as-code for your service under the `copilot/` directory.
```
copilot
└── vote
    └── manifest.yml
```
The `manifest.yml` file under `vote/` holds the common configuration for a "load balanced web service" pattern.
For example, you can setup configuration for your ECS task size, exposed port, environment variables,
secrets from AWS Secrets Manager or AWS Systems Manager Parameter Store.

## Deleting the service
If you'd like to delete only the service from the "voting-app" application. 
```bash
$ copilot svc delete
```
If you'd like to delete the entire application including other services and deployment environments:
```bash
$ copilot app delete
```
# matmult

In this project, 04 AWS lambda functions were built in which they were developed respectively in pure Go, pure Python, Go using the GoNum library and Python with the NumPy library. The functions present the matrix multiplication algorithm. The four functions were part of an experiment carried out, in which the objective was to evaluate the performance of the Go and Python languages with high-performance applications in a serverless paradigm on the AWS platform. In addition, it has a model for implementing the architecture in SAM for application automation, which includes the construction of the 04 functions, API Gateway definitions and construction of the table in DynamoDB.

This is a sample template for matmult - Below is a brief explanation of what we have generated for you:

```bash
.
├── Makefile                  <-- Make to automate build
├── README.md                   <-- This instructions file
├── gonum                    <-- Source code for a lambda function in Go with GoNum for matrix multiplication
│   ├── main.go                       <-- Lambda function code
├── multgo                    <-- Source code for a lambda function in Go for matrix multiplication
│   ├── main.go                       <-- Lambda function code
│   └── main_test.go                <-- Unit tests
├── multpython             <-- Source code for a lambda function in Python for matrix multiplication
│   ├── main.py                       <-- Lambda function code
├── numpy                   <-- Source code for a lambda function in Python with NumPy for matrix multiplication
│   ├── main.py                       <-- Lambda function code
├── result                     <-- Source code for a lambda function in Go to retrieve result
│   ├── main.go                       <-- Lambda function code
└── template.yaml        <-- Infrastructure model with SAM
```

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Setup process

### Installing dependencies & building the target 

In this example we use the built-in `sam build` to automatically download all the dependencies and package our build target.   
Read more about [SAM Build here](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-build.html) 

The `sam build` command is wrapped inside of the `Makefile`. To execute this simply run
 
```shell
make
```

### Local development

**Invoking function locally through local API Gateway**

```bash
sam local start-api
```

If the previous command ran successfully you should now be able to hit the following local endpoint to invoke your function `http://localhost:3000/hello`

** SAM CLI ** is used to emulate Lambda and API Gateway locally and uses our `template.yaml` to understand how to initialize this environment (runtime, where the source code is, etc.) - The following excerpt is an example what the CLI will read to initialize an API and its routes:

```yaml
...
Events:
    HelloWorld:
        Type: Api # More info about API Event Source: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#api
        Properties:
            Path: /hello
            Method: get
```

## Packaging and deployment

The AWS Lambda runtime for different languages requires a simple folder with the executable generated in the compilation step. SAM will use the `CodeUri` property to know where to look for the application. For example:

```yaml
...
    FirstFunction:
        Type: AWS::Serverless::Function
        Properties:
            CodeUri: hello_world/
            ...
```

To deploy your application for the first time, run the following in your shell:

```bash
sam deploy --guided
```

The command will package and deploy your application to AWS, with a series of prompts:

* **Stack Name**: The name of the stack to deploy to CloudFormation. This should be unique to your account and region, and a good starting point would be something matching your project name.
* **AWS Region**: The AWS region you want to deploy your app to.
* **Confirm changes before deploy**: If set to yes, any change sets will be shown to you before execution for manual review. If set to no, the AWS SAM CLI will automatically deploy application changes.
* **Allow SAM CLI IAM role creation**: Many AWS SAM templates, including this example, create AWS IAM roles required for the AWS Lambda function(s) included to access AWS services. By default, these are scoped down to minimum required permissions. To deploy an AWS CloudFormation stack which creates or modified IAM roles, the `CAPABILITY_IAM` value for `capabilities` must be provided. If permission isn't provided through this prompt, to deploy this example you must explicitly pass `--capabilities CAPABILITY_IAM` to the `sam deploy` command.
* **Save arguments to samconfig.toml**: If set to yes, your choices will be saved to a configuration file inside the project, so that in the future you can just re-run `sam deploy` without parameters to deploy changes to your application.

You can find your API Gateway Endpoint URL in the output values displayed after deployment.
```
### Testing

To test the various functions, we can make a POST request for the endpoint of the function you want, you must send the order (Matrix) and programming language (Lambda Function) in the request.

**** Test by API Gateway (json), for example:
{"id": "48", "order": 800, "language": "go", "time": "", "time": ""}

**** Curl test, example:
curl -H "Content-Type: application / json" \
         -X POST \
         -d '{"id": "48", "order": 800, "language": "numpy", "time": "", "time": ""}' \
         $ {Invocation URL} / function-endpoint
```
# Appendix

### Golang installation

Please ensure Go 1.x (where 'x' is the latest version) is installed as per the instructions on the official golang website: https://golang.org/doc/install

A quickstart way would be to use Homebrew, chocolatey or your linux package manager.

#### Homebrew (Mac)

Issue the following command from the terminal:

```shell
brew install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
brew update
brew upgrade golang
```

#### Chocolatey (Windows)

Issue the following command from the powershell:

```shell
choco install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
choco upgrade golang
```

## Bringing to the next level

Here are a few ideas that you can use to get more acquainted as to how this overall process works:

* Create an additional API resource (e.g. /hello/{proxy+}) and return the name requested through this new path
* Update unit test to capture that
* Package & Deploy

Next, you can use the following resources to know more about beyond hello world samples and how others structure their Serverless applications:

* [AWS Serverless Application Repository](https://aws.amazon.com/serverless/serverlessrepo/)

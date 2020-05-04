package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)


type AWSCloudResource struct {
}

type printAWSR func(string) string

type AWSCloud struct {
	region 	string
	limit   string
	Bucket	[]Bucket
	Lambda []Lambda
	DynamoDB []DynamoDB
}

func printSubHeader(r string) string {
	return  "AWS" + r
}

func (a *AWSCloudResource) printHeader(fn printAWSR, r string) {
	fmt.Println("************************************")
	fmt.Printf("* %s \n", fn(r))
}

func OnError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}


func (a *AWSCloudResource) handleError(x string, err error) {
	if err != nil {
		OnError("%s, %v", x, err)
	}
}

func (a *AWSCloud) CreateSession() interface{}{
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sess
}

func (a *AWSCloudResource) ListStorageBuckets(s interface{}) []Bucket {
	sess := s.(*session.Session)
	svc := s3.New(sess)
	result, err := svc.ListBuckets(nil)
	a.handleError("Unable to list buckets", err)
	a.printHeader(printSubHeader, "S3 Buckets:")

	var buckets []Bucket
	for _, b := range result.Buckets {
		n := Bucket{BucketName: aws.StringValue(b.Name), CreationDate: aws.TimeValue(b.CreationDate).String()}
		buckets = append( buckets, n)

		fmt.Printf("* %s created on %s\n",aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
	return buckets
}

func (a *AWSCloudResource) ListLambda(s interface{}) []Lambda {

	sess := s.(*session.Session)
	lambdaSvc := lambda.New(sess)
	result, err := lambdaSvc.ListFunctions(nil)
	a.handleError("Unable to list Lambda", err)
	a.printHeader(printSubHeader, "Lambdas:")

	var lambda []Lambda
	for _, b := range result.Functions {
		n := Lambda{FunctionName: aws.StringValue(b.FunctionName), FunctionArn: aws.StringValue(b.FunctionArn)}
		lambda = append(lambda, n)
		fmt.Printf("* %s :: function arn ::  %s\n",
			aws.StringValue(b.FunctionName), aws.StringValue(b.FunctionArn))
	}

	return lambda
}

func (a *AWSCloudResource) ListDBTables(s interface{}) []DynamoDB {
	sess := s.(*session.Session)
	dynamoDBSvc := dynamodb.New(sess)
	result, err := dynamoDBSvc.ListTables(nil)
	a.handleError("Unable to list DynamoDB tables", err)
	a.printHeader(printSubHeader, "DynamoDB Tables:")

	var tables []DynamoDB

	for _, b := range result.TableNames {
		n := DynamoDB{ DatabaseName: aws.StringValue(b)}
		tables = append(tables, n)
		fmt.Printf("* %s \n",
			aws.StringValue(b))
	}
	return tables
}

func (a *AWSCloudResource) ListApiGatewayEndpoints(s interface{}) {
	sess := s.(*session.Session)
	gatewaySvc := apigateway.New(sess)
	gatewayInput := &apigateway.GetRestApisInput{
		Limit: aws.Int64(5),
	}
	result, _ := gatewaySvc.GetRestApis(gatewayInput)
	a.printHeader(printSubHeader, "APIGateway Endpoints")
	for _, b := range result.Items {
		fmt.Printf("* %s ::description:  %s\n",
			aws.StringValue(b.Name), aws.StringValue(b.Description))
	}
}

func (a *AWSCloudResource) ListCloudFormationStackSets(s interface{}){
	sess := s.(*session.Session)
	cfSvc := cloudformation.New(sess)

	result, err := cfSvc.ListStackSets(nil)
	a.handleError("Unable to list StackSets", err)
	a.printHeader(printSubHeader, "StackSets:")
	for _, b := range result.Summaries {
		fmt.Printf("* %s %s\n",
			aws.StringValue(b.StackSetName), aws.StringValue(b.Description))
	}
}

func (a *AWSCloudResource) ListCloudFormationStack(s interface{}){
	sess := s.(*session.Session)
	cfSvc := cloudformation.New(sess)

	result, err := cfSvc.ListStacks(nil)
	a.handleError("Unable to list StackSets", err)
	a.printHeader(printSubHeader, "Stacks:")
	for _, b := range result.StackSummaries {
		fmt.Printf("* %s %s\n",
			aws.StringValue(b.StackName), aws.StringValue(b.StackId))
	}
}

func (a *AWSCloudResource) ListCloudFront(s interface{}){
	sess := s.(*session.Session)
	cfSvc := cloudfront.New(sess)

	result, err := cfSvc.ListDistributions(nil)
	a.handleError("Unable to list CloudFront Dists", err)
	a.printHeader(printSubHeader, "CloudFront Distribution:")
	for _, b := range result.DistributionList.Items {
		fmt.Printf("* %s %s\n",
			aws.StringValue(b.Id), aws.StringValue(b.ARN))
	}
}

func (a *AWSCloudResource) ListMetrics(s interface{}){
	sess := s.(*session.Session)
	cwSvc := cloudwatch.New(sess)

	result, err := cwSvc.ListMetrics(nil)
	a.handleError("Unable to list CloudWatch Metrics", err)
	a.printHeader(printSubHeader, "CloudWatch Metrics:")
	for _, b := range result.Metrics {
		fmt.Printf("* %s %s\n",
			aws.StringValue(b.MetricName), aws.StringValue(b.Namespace))
	}
}

func (a * AWSCloud) ListPricingDetails () {
	fmt.Print("Pricing .....")
}

func (a *AWSCloud) GetResources()  {
	awsR := AWSCloudResource{}

	sess := a.CreateSession()
	a.Bucket = awsR.ListStorageBuckets(sess)
	a.Lambda = awsR.ListLambda(sess)
	a.DynamoDB = awsR.ListDBTables(sess)
	awsR.ListApiGatewayEndpoints(sess)
	awsR.ListCloudFormationStack(sess)
	awsR.ListCloudFront(sess)
	awsR.ListCloudFormationStackSets(sess)

}
package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
	"log"
)

func main() {
	ctx := context.Background()

	opt := option.WithCredentialsFile("./data/piapp-ios-dev-firebase-adminsdk-jnh9b-308df48c0b.json")
	//config := &firebase.Config{ProjectID: "my-project-id"}
	//opt := option.WithCredentialsJSON(serviceAccountKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("error client : %v \n", err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "GEkn0b5w:APA91bEwSdi2bxVhrZBqSqiRVA7EkUWvqGI0l4nk6RPQJN4Irq3HZZ0IpCd9BB0wE7avJpTcS0jai1SioClYr6u7Cin9U5d2J3BNBxF6iRDBJOH4iiMUojJdp1ytGkJzUANohNLG0xaf"

	topic := "zzzz"
	//
	//// [START subscribe_golang]
	// These registration tokens come from the client FCM SDKs.
	registrationTokens := []string{}
	for i := 0; i < 1000; i++ {
		registrationTokens = append(registrationTokens, registrationToken)
	}

	// Subscribe the devices corresponding to the registration tokens to the
	// topic.
	response, err := client.SubscribeToTopic(ctx, registrationTokens, topic)
	if err != nil {
		log.Fatalln(err)
	}
	// See the TopicManagementResponse reference documentation
	// for the contents of response.
	fmt.Println(response.SuccessCount, "tokens were subscribed successfully")

	//subEnd := time.Now()

	// [START send_to_topic_golang]
	// The topic name can be optionally prefixed with "/topics/".
	//uid := "25696773511053390"

	badge := int(48)

	var aps = &messaging.Aps{
		AlertString:      "\\U5594\\U5594\\U4eba\\U554a\\U4ed8\\U6b3e$100\\U7d66\\U60a8\\Uff0c\\U8acb\\U78ba\\U8a8d\\U662f\\U5426\\U6536\\U4e0b\\U3002",
		Badge:            &badge,
		ContentAvailable: true,
		Sound:            "bingbong.aiff",
	}

	var Payload = &messaging.APNSPayload{
		Aps: aps,
	}

	var Apns = &messaging.APNSConfig{
		Payload: Payload,
	}

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"client": "CL402FL0MAOITBDK",
			"data":   "{\\\"result\\\":\\\"success\\\",\\\"error_code\\\":\\\"0000\\\",\\\"error_message\\\":\\\"\\\",\\\"payment_amount\\\":\\\"$111\\\",\\\"title\\\":\\\"\\\\u4ed8\\\\u6b3e\\\\u6210\\\\u529f\\\",\\\"description\\\":\\\"\\\\u5e97\\\\u5bb6\\\\u540d\\\\u7a31 / \\\\u5168\\\\u5bb6\\\\u4fbf\\\\u5229\\\\u5546\\\\u5e97\\\\n\\\\u4ed8\\\\u6b3e\\\\u65b9\\\\u5f0f /  6080\\\\n\\\\n\\\\u203b P \\\\u5e63\\\\u56de\\\\u994b\\\\u4ee5 P \\\\u5e63\\\\u5e33\\\\u6236\\\\u70ba\\\\u4e3b\\\"}",
			"title":  "Pi \\U62cd\\U9322\\U5305 \\U4ed8\\U6b3e\\U6210\\U529f",
			"type":   "type",
		},
		APNS:  Apns,
		Token: registrationToken,
	}

	// Send a message to the devices subscribed to the provided topic.
	responses, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", responses)

}

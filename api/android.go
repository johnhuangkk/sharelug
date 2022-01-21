package api

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
	"log"
)

func main() {
	ctx := context.Background()

	opt := option.WithCredentialsFile("./data/piapp-eada5-firebase-adminsdk-uh64n-627b5bd507.json")
	//config := &firebase.Config{ProjectID: "my-project-id"}
	//opt := option.WithCredentialsJSON(serviceAccountKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		//fmt.Println(err)
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		//fmt.Println(err)
		log.Fatalf("error client : %v \n", err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "cR195bmrsws:APA91bFpPDP2HlpYi20x2SloMb9Mwyb3_Eci1ehWuplDAuX3dZIQshbcYBYVajJa2OdmJlku9CSTmCO7yMgLbSS1CWDhDllizfDeQrtehDbQiINb_7jQ0vnolMBP4y3VzAgXj0CFDP58"

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

	//badge := int(48)
	//
	//var aps = &messaging.Aps{
	//	AlertString: "\\U5594\\U5594\\U4eba\\U554a\\U4ed8\\U6b3e$100\\U7d66\\U60a8\\Uff0c\\U8acb\\U78ba\\U8a8d\\U662f\\U5426\\U6536\\U4e0b\\U3002",
	//	Badge: &badge,
	//	ContentAvailable:true,
	//	Sound:"bingbong.aiff",
	//}
	//
	//var Payload = &messaging.APNSPayload{
	//	Aps: aps,
	//}
	//
	//var Apns = &messaging.APNSConfig{
	//	Payload: Payload,
	//}

	var Notification = &messaging.Notification{
		Body: "Android 測試測試測試",
	}

	datas := map[string]string{
		"client":      "CL402FL0MAOITBDK",
		"type":        "MESSAGE_CREATED",
		"title":       "\\U5594\\U5594\\U4eba\\U554a",
		"description": "",
		"data":        "{\"id\":\"TM392GGIWJW0PI6G\",\"room_id\":\"TR672GGIWJUNB0BW\",\"user_id\":\"UP063CJHMIFQ6688\",\"listener_user_id\":\"UP883CFJSQLHRHYG\",\"transaction_id\":\"TX2032342328163848\",\"read_user_ids\":\"[]\",\"read_count\":0,\"type\":\"2\",\"content\":\"\",\"status\":\"1\",\"metadata\":\"\",\"updated_at\":\"2020-05-14 18:23:29\",\"created_at\":\"2020-05-14 18:23:29\"}",
	}

	data, err := json.Marshal(datas)
	if err != nil {
		panic(err)
	}

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"message": string(data),
		},
		//APNS: Apns,
		Notification: Notification,
		Token:        registrationToken,
	}

	// Send a message to the devices subscribed to the provided topic.
	responses, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", responses)

}

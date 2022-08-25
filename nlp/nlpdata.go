package nlp

const nlpOriginalData = `{
    "intents": [
        {
            "tag": "retrain",
            "patterns": [
                "retrain your nlp data",
                "retrain your natural language data",
                "upgrade the training data"
            ]
        },
        {
            "tag": "greet",
            "patterns": [
                "Hello",
                "Good morning",
                "good afternoon",
                "Hi",
                "Howdy",
                "Greetings"
            ]
        },
        {
            "tag": "goodbye",
            "patterns": [
                "goodbye",
                "see you",
                "farewell",
                "good night"
            ]
        },
        {
            "tag": "current weather",
            "patterns": [
                "What is the weather like",
                "What is the weather now",
                "what is the weather"
            ]
        },
        {
            "tag": "weather future",
            "patterns": [
                "what will the weather be like tomorrow",
                "what will the weather be like in {{number}} days",
                "what will the weather be like on {{day}}",
                 "what is tomorrows weather"
            ]
        },
        {
            "tag": "weather past",
            "patterns": [
                "what was the weather like yesterday",
                "what was the weather like {{number}} days ago",
                "what was the weather like on {{day}}",
                "what was yesterdays weather"
            ]
        },
        {
            "tag": "who",
            "patterns": [
                "who is {{subject}}",
                "who was {{subject}}"
            ]
        },
        {
            "tag": "what",
            "patterns": [
                "what is {{subject}}",
                "what was {{subject}}"
            ]
        },
        {
            "tag":"who_are_you",
            "patterns":["who are you", "what is your name","what are you"]
        }, {
            "tag":"how_are_you",
            "patterns":["how are you", "how do you feel","are you ok", "do you feel fine","do you feel ok", "are you fine","how are you feeling"]
        }
    ]
}`

// "relearn your nlp data",
// "relearn your natural language data",
// "learn your nlp data",
// "learn your natural language data",
// "read your nlp data",
// "read your natural language data",
// "update your nlp data",
// "update your natural language data",
// "upgrade your natural language data",
// "upgrade your nlp data",
// "retrain your training data",
// "relearn your training data",
// "learn your training data",
// "read your training data",
// "update your training data",
// "upgrade your training data",
// "retrain the nlp data",
// "retrain the natural language data",
// "relearn the nlp data",
// "relearn the natural language data",
// "learn the nlp data",
// "learn the natural language data",
// "read the nlp data",
// "read the natural language data",
// "update the nlp data",
// "update the natural language data",
// "upgrade the natural language data",
// "upgrade the nlp data",
// "retrain the training data",
// "relearn the training data",
// "learn the training data",
// "read the training data",
// "update the training data",

// {
//     "tag": "math_add",
//     "patterns": [
//         "add {{num}} and {{num}}",
//         "add {{num}} to {{num}}",
//         "what is the sum of {{num}} and {{num}}",
//         "what is {{num}} add {{num}}",
//         "what is {{num}} plus {{num}}",
//         "sum {{num}} and {{num}}"]
// },{

// "tag":"math_subtract",
// "patterns":[
//         "take {{num}} from {{num}}",
//         "subtract {{num}} from {{num}}",
//         "what is {{num}} subtracted from {{num}}",
//         "what is {{num}} minus {{num}}"]
// },{
//     "tag":"math_multiply",
//     "patterns":[
//         "multiply {{num}} and {{num}}",
//         "multiply {{num}} to {{num}}",
//         "product of {{num}} and {{num}}",
//         "what is {{num}} multiplied by {{num}}",
//         "what is {{num}} times {{num}}",
//         "what is the product of {{num}} and {{num}}"]
// },{
//     "tag":"math_division",
//     "patterns":[
//         "divide {{num}} by {{num}}",
//         "divide {{num}} into {{num}}",
//         "what is {{num}} divided by {{num}}",
//         "what is {{num}} shared by {{num}}"]
// },{
//     "tag":"math_modulus",
//     "patterns":[
//         "remainder from {{num}} divided by {{num}}",
//         "what is left if we divide {{num}} by {{num}}"
//     ]
// },

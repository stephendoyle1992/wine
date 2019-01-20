// See https://github.com/dialogflow/dialogflow-fulfillment-nodejs
// for Dialogflow fulfillment library docs, samples, and to report issues
'use strict';
 
const functions = require('firebase-functions');
const {WebhookClient} = require('dialogflow-fulfillment');
const {Card, Suggestion} = require('dialogflow-fulfillment');
 
process.env.DEBUG = 'dialogflow:debug'; // enables lib debugging statements

let loadedWines = []
let curSelected = 0


 
exports.dialogflowFirebaseFulfillment = functions.https.onRequest((request, response) => {
  const agent = new WebhookClient({ request, response });
  
 
  function welcome(agent) {
    var request = require('request');
    
    
    return new Promise((resolve, reject) => {
        let country = agent.parameters.country;
        let type = agent.parameters.red_white;
        let price = agent.parameters.price;
        let sorting = agent.parameters.sort_wine;
        loadedWines = []
        
        agent.add(`Ok so you want a ${type} wine from ${country}, up to ${price} dollars? Let's see.`);
        
        if (country === "United States of America") {
            country = "US";
        }
        
        let url = `https://immense-escarpment-78210.herokuapp.com/api/wines/?type=${type}&price=${price}&country=${country}&sorting=${sorting}`
        console.log(url)
        request.get(url, (error, response, body) => {
            let data = JSON.parse(body);
            let output;
            console.log(data);
            console.log(data.length);
            if (data.length == 0) {
                output = agent.add(`I didn't find shit.`);
            } else {
                loadedWines = data;
                output = agent.add(`I'd recommend the ${loadedWines[0].title.String}. Twitter user ${loadedWines[0].tasterTwitterHandle.String} 
                    described it as ${loadedWines[0].description.String}.
                    The price was last seen at ${loadedWines[0].price.Int64} dollars american.`);
                    curSelected = 1;
            }
      resolve(output);
    });
  });
  }
  
  function readNext(agent) {
      console.log(loadedWines)
      if (curSelected < loadedWines.length) {
          agent.add(`I'd also recommend the ${loadedWines[curSelected].title.String}. Twitter user ${loadedWines[curSelected].tasterTwitterHandle.String} 
                    described it as ${loadedWines[curSelected].description.String}.
                    The price was last seen at ${loadedWines[curSelected].price.Int64} dollars american.`)
                    curSelected++;
      } else {
          agent.add('I dont have any more suggestions')
      }
  }
 
  function fallback(agent) {
    agent.add(`I didn't understand`);
    agent.add(`I'm sorry, can you try again?`);
  }

  let intentMap = new Map();
  intentMap.set('Default Welcome Intent', welcome);
  intentMap.set('Default Fallback Intent', fallback);
  intentMap.set(`next suggestion`, readNext);

  agent.handleRequest(intentMap);
});

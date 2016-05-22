var request = require('request')

var options = {
	url: "https://www.reddit.com/r/dota2/hot.json",
	headers: {
		'User-Agent': 'dota2logger by /u/VRCkid'
	}
};

request(options, function(error, response, body) {
	obj = JSON.parse(body);
	
	obj.data.children.forEach(function(element) {
		console.log(element.data.title + " " + element.data);
	});
});

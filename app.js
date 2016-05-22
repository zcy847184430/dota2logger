var mongoose = require('mongoose');
var request = require('request')

var postSchema = new mongoose.Schema({
	title: String,
	permalink: String,
	url: String,
	self: Boolean,
	stickied: Boolean,
	score: Number,
	flair: String,
	author: String,
	date: {type: Date, default: Date.now() }
});

var Post = mongoose.model('Post', postSchema);

function main() {
	console.log("Started main");
	
	var options = {
		url: "https://www.reddit.com/r/dota2/hot.json",
		headers: {
			'User-Agent': 'dota2logger by /u/VRCkid'
		}
	};

	request(options, function(error, response, body) {
		obj = JSON.parse(body);
		var counter = 0;
		
		obj.data.children.forEach(function(element) {

			if(counter < 16) {
				counter++;
				console.log(element.data.title + " " + element.data);

				Post.findOne({permalink: element.data.permalink}, (err, post) => {
					if(post !== null) {
						console.log("Updating post score");

						post.score = element.data.score;

						post.save( (err, post) => {
							if(err) throw err;
						});
					}
					else {
						console.log("New post!");

						var time = new Date(0);
						time.setUTCSeconds(element.data.created_utc);

						var post = new Post({
							title: element.data.title,
							permalink: element.data.permalink,
							url: element.data.url,
							self: element.data.is_self,
							stickied: element.data.stickied,
							score: element.data.score,
							flair: element.data.link_flair_text,
							author: element.data.author,
							date: time							
						});

						post.save( (err, post) => {
							if(err) throw err;
						});
					}
				});
			}

		});
	});

}

mongoose.connect('mongodb://vatyx:vatyx@ds041613.mlab.com:41613/dota2logger')

var db = mongoose.connection

console.log("Hello");

function what() {
	console.log("yep");
}

db.on('error', () => { console.log("welp"); })
db.once('open', () => { setInterval(main, 900000) })




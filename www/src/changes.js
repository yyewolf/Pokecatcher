function autochange(){
	autocatcher = document.getElementById("catcher").checked;
	ws.send('{"action":"aca","state":'+autocatcher+'}');
	if(autocatcher){
		NotifTitle = 'I will now auto catch Pokémon !';
		notify('info', NotifTitle);
	}else{
		NotifTitle = 'I won\'t auto catch Pokémon !';
		notify('info', NotifTitle);
	};
};

function alchange(){
	autoleveling = document.getElementById("al").checked;
	ws.send('{"action":"al","state":'+autoleveling+'}');
	if(autoleveling){
		NotifTitle = 'I will now auto level your Pokémon !';
		notify('info', NotifTitle);
	}else{
		NotifTitle = 'I won\'t auto level your Pokémon !';
		notify('info', NotifTitle);
	};
};

function duplichange(){
	duplicate = document.getElementById("duplicate").checked;
	ws.send('{"action":"duplicate","state":'+duplicate+'}');
	if(duplicate){
		NotifTitle = 'I will now ignore duplicate Pokémon!';
		notify('info', NotifTitle);
	}else{
		NotifTitle = 'I woill now catch every Pokémon !';
		notify('info', NotifTitle);
	};
};

function filterchange(){
	filter = document.getElementById("filter").checked;
	ws.send('{"action":"filter","state":'+filter+'}');
	if(filter){
		NotifTitle = 'I will now filter every Pokémon (using pre-made filter)!';
		notify('info', NotifTitle);
	}else{
		NotifTitle = 'I won\'t filter every Pokémon (using pre-made filter)!';
		notify('info', NotifTitle);
	};
};

function customfilterchange(){
	customfilters = document.getElementById("customfilter").checked;
	ws.send('{"action":"customfilters","state":'+customfilters+'}');
	if(customfilters){
		NotifTitle = 'I will now filter every Pokémon (using custom filters)!';
		notify('info', NotifTitle);
	}else{
		NotifTitle = 'I won\'t filter every Pokémon (using custom filters)!';
		notify('info', NotifTitle);
	};
};

function alertschange(){
	alerts = document.getElementById("alerts").checked;
	ws.send('{"action":"alerts","state":'+alerts+'}');
	if(alerts){
		NotifTitle = 'I will now alert if I can!';
		notify('info', NotifTitle);
	}else{
		NotifTitle = 'I won\'t alert anymore!';
		notify('info', NotifTitle);
	};
};


function aliaschange(){
	aliases = document.getElementById("aliases").checked;
	ws.send('{"action":"aliases","state":'+aliases+'}');
	if(aliases){
		NotifTitle = 'I will now catch every Pokémon with aliases !';
		notify('info', NotifTitle);
	}else{
		NotifTitle = 'I won\'t catch Pokémon with aliases !';
		notify('info', NotifTitle);
	};
};

function refresh(){
	ws.send('{"action":"refresh"}');
	notify('info', 'Refreshing your pokemon list...');
};

function refreshmoves(){
	ws.send('{"action":"refreshmoves"}');
	NotifTitle = 'Refreshing your pokemon moves list...';
	notify('info', NotifTitle);
};

function catchpokemon(name, id, command){
	ws.send('{"action":"catch","name":"'+name+'","channelid":"'+id+'","command":"'+command+'"}');
};

function setmove(move, pokemonname, number, channel){
	ws.send('{"action":"learn","pokemonname":"'+pokemonname+'","movenumber":'+number+',"movename":"'+move+'","channelset":"'+channel+'"}');
	NotifTitle = 'Learned move to ' + pokemonname.replace(/\**/g, '').replace('**', '').replace(/ /g, '').replace(/ /g, '').replace(/é/g, 'e').replace(/è/g, 'e').replace('♂', '').replace('♀', '') + ' !';
	NotifIcon = './img/' + pokemonname.replace(/\**/g, '').replace('**', '').replace(/ /g, '').replace(/ /g, '').replace(/é/g, 'e').replace(/è/g, 'e').replace('♂', '').replace('♀', '') + '.png';
	notify('info', NotifTitle, NotifIcon);
};

function spamchange(){
	var spamwillbeactive = document.getElementById("spammer").checked;
	if(spamwillbeactive == true){
		var spaminterval = document.getElementById("spaminterval").value;
		var spammessage = document.getElementById("spamtext").value;
		if(spammessage != undefined && spaminterval != undefined) {
			if(spammessage != " ") {
				if(spammessage != "") {
					console.log(spammessage);
					ws.send('{"action":"spam","state":true,"message":"'+spammessage+'","interval":'+spaminterval+'}');
					NotifTitle = 'Spamming enabled !';
					notify('info', NotifTitle);
				}else{
					alert('You didn\'t set any message/interval to spam !')
					document.getElementById("spammer").checked = false;
				};
			}else{
				alert('You didn\'t set any message/interval to spam !')
				document.getElementById("spammer").checked = false;
			};
		}else{
			alert('You didn\'t set any message/interval to spam !')
			document.getElementById("spammer").checked = false;
		};
	}else{
		ws.send('{"action":"spam","state":false}');
		NotifTitle = 'Spamming disabled !';
		notify('info', NotifTitle);
	};
};

function dothis(action,number,name){
	switch(action) {
		case 'release' :
			NotifTitle = 'Releasing :';
			NotifIcon = './img/' + name + '.png';
			NotifMessage = name;
			notify('info', NotifTitle, NotifIcon, NotifMessage);
			ws.send('{"action":"release","number":'+number+'}');
			break;
		case 'remove' :
			NotifTitle = 'Removing from the list :';
			NotifIcon = './img/' + name + '.png';
			NotifMessage = name;
			notify('info', NotifTitle, NotifIcon, NotifMessage);
			ws.send('{"action":"remove","number":'+number+'}');
			break;
		case 'select' :
			NotifTitle = 'Selecting :';
			NotifIcon = './img/' + name + '.png';
			NotifMessage = name;
			notify('info', NotifTitle, NotifIcon, NotifMessage);
			ws.send('{"action":"select","number":'+number+',"name":"'+name+'"}');
			break;
		case 'addfav' :
			NotifTitle = 'Adding Favorite :';
			NotifIcon = './img/' + name + '.png';
			NotifMessage = name;
			notify('info', NotifTitle, NotifIcon, NotifMessage);
			ws.send('{"action":"addfav","number":'+number+'}');
			break;
		case 'remfav' :
			NotifTitle = 'Removing Favorite :';
			NotifIcon = './img/' + name + '.png';
			NotifMessage = name;
			notify('info', NotifTitle, NotifIcon, NotifMessage);
			ws.send('{"action":"remfav","number":'+number+'}');
			break;
		case 'nickname' :
			var nickname = document.getElementById("nicknametext").value;
			if(nickname != undefined) {
				if(nickname != " ") {
					if(nickname != "") {
						ws.send('{"action":"nickname","nickname":"'+nickname+'"}');
					};
				};
			};
			break;
		case 'tokenchange' :
			var tokenchange = document.getElementById("tokentext").value;
			if(tokenchange != undefined) {
				if(tokenchange != " ") {
					if(tokenchange != "") {
						NotifTitle = 'Restart the bot to apply the new token !';
						notify('danger', NotifTitle);
						ws.send('{"action":"tokenchange","token":"'+tokenchange+'"}');
					};
				};
			};
			break;
		case 'autodelay' :
			var delay = document.getElementById("delayms").value;
			if(delay != undefined) {
				if(delay != " ") {
					if(delay != "") {
						NotifTitle = 'The delay has been changed !';
						notify('danger', NotifTitle);
						ws.send('{"action":"autodelaychange","delay":'+delay+'}');
					};
				};
			};
			break;
		case 'queue' :
			var queuelist = document.getElementById("alt").value;
			if(queuelist != undefined) {
				if(queuelist != " ") {
					if(queuelist != "") {
						NotifTitle = 'The priority queue has been changed !';
						notify('danger', NotifTitle);
						ws.send('{"action":"queuelist","change":"'+queuelist+'"}');
					};
				};
			};
			break;
	};
};

function refreshwhitelist(id){
	ws.send('{"action":"whitelist","serverid":"'+id+'","serverstate":'+document.getElementById("server"+id).checked+'}');
};

function refreshpkmnwhitelist(name, checked){
	pokewhitelist[name] = checked
	ws.send('{"action":"pokemonwhitelist","name":"'+name+'","state":'+checked+'}');
};

function importpkmnwhitelist(){
	obj = {
		action:"importwhitelist",
		change:JSON.stringify(pokewhitelist)
	}
	ws.send(JSON.stringify(obj));
	NotifTitle = 'Imported a new whitelist !';
	notify('info', NotifTitle);
};

function changeprefix(type,prefix){
	NotifTitle = 'This prefix has been changed!';
	notify('info', NotifTitle);
	ws.send('{"action":"prefixchange","type":"'+type+'","prefix":"'+prefix+'"}');
};

function changeautoleveler(max){
	NotifTitle = 'The AutoLeveler max level has been changed!';
	notify('info', NotifTitle);
	ws.send('{"action":"alchange","change":"'+max+'"}');
};

function SaveFilters(){
	NotifTitle = 'Your custom filters have been changed!';
	notify('info', NotifTitle);
	ws.send('{"action":"filterschange","filters":'+JSON.stringify(filters)+'}');
};
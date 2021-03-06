if ("WebSocket" in window) {
	var ws = new WebSocket("ws://localhost:"+websocket);
	ws.onopen = function() {};
	ws.onerror = function (err) { 
		console.log(err);
	}
	ws.onmessage = function (evt) { 
		var message = evt.data;
		var msg = JSON.parse(message);
		switch(msg["action"]) {
			case 'selected' :
				document.getElementById("img").innerHTML = '<img src="./img/'+ClearAll(msg['name'])+'.png" class="img-circle" width="50" height="50">';
				break;
			case 'registered' :
				document.getElementById("alertbox").parentNode.removeChild(document.getElementById("alertbox"));
				break;
			case 'refreshlist' :
				if(msg['listobj'] == ""){
					LoadList("");
				}
				listobj = JSON.parse(msg['listobj']);
				LoadList(listobj);
				break;
			case 'removefromlist' :
				if(msg['listobj'] == ""){
					LoadList("");
				}
				listobj = JSON.parse(msg['listobj']);
				LoadList(listobj);
				break;
			case 'addpokemon' :
				if(msg['listobj'] == ""){
					LoadList("");
				}
				listobj = JSON.parse(msg['listobj']);
				LoadList(listobj);
				break;
			case 'notification' :
				if(autocatcher == false){
					NotifIcon = './img/' + ClearAll(msg['name']) + '.png';
					NotifTitle = 'A ' + ClearAll(msg['name']) + ' has spawned !';
					NotifMessage = 'Type ' + msg['command'] + ' ' + ClearAll(msg['name']) + ' to catch it ! <a onclick="catchpokemon(\'' + msg['name'].replace('é','e') + '\',\'' + msg['channelid'] +'\', \'' + msg['command'] +'\')" href="#">Catch it!</a>';
					NotifChannelOf = 'In #'+msg['channelname']+' of '+msg['server']+'</span>';
					notify('info', NotifTitle, NotifIcon, NotifMessage, 3000, NotifChannelOf);
				}else{
					NotifIcon = './img/' + ClearAll(msg['name']) + '.png';
					NotifTitle = 'A ' + ClearAll(msg['name']) + ' has spawned !';
					NotifMessage = 'Type ' + msg['command'] + ' ' + ClearAll(msg['name']) + ' to catch it !';
					NotifChannelOf = 'In #'+msg['channelname']+' of '+msg['server']+'</span>';
					notify('info', NotifTitle, NotifIcon, NotifMessage, 3000, NotifChannelOf);
				};
				break;
			case 'warn' :
				switch(msg["message"]) {
					case 'couldnt-normal':
						NotifIcon = './img/' + ClearAll(msg['name']) + '.png';
						NotifTitle = 'I couldn\'t catch ' + ClearAll(msg['name']) + ' !';
						NotifChannelOf = 'In #'+msg['channelname']+' of '+msg['server']+'</span>';
						notify('warning', NotifTitle, NotifIcon, 0, 3000, NotifChannelOf);
						break;
					case 'couldnt-duplicate':
						NotifIcon = './img/' + ClearAll(msg['name']) + '.png';
						NotifTitle = 'I didn\'t catch' + ClearAll(msg['name']) + ' because you have it!';
						NotifChannelOf = 'In #'+msg['channelname']+' of '+msg['server']+'</span>';
						notify('warning', NotifTitle, NotifIcon, 0, 3000, NotifChannelOf);
						break;
					case 'nickname-succes':
						NotifIcon = './src/check.png';
						NotifTitle = 'Successfully changed the nickname of the current pokemon !';
						notify('danger', NotifTitle, NotifIcon, 0, 3000);
						break;
					case 'nickname-failed':
						NotifTitle = 'Couldn\'t change the nickname of the current pokemon !';
						notify('danger', NotifTitle, 0, 0, 3000);
						break;
				}
				break;
			case 'refreshmovelist' :
				NotifTitle = 'Your pokemon moves list has been refreshed';
				notify('info', NotifTitle, 0, 0, 4000);
				pokemonname = msg['pokemonname'];
				channelset = msg['channelmovesid'];
				movelist = msg['moves'].split(';');
				document.getElementById("moveslist").innerHTML = "";
				fetch('./attacks.json').then(res => res.json())
					.then((lines) => {
						console.log(lines)
						for(var k=0 ; k < movelist.length ; k++){
							console.log(lines[movelist[k]])
							if(lines[movelist[k]] != undefined) {
								var html = '';
								html += '<p style="display:inline-block;">';
									html += '<form style="float: left; padding: 5px;" id="formmove">';
										html += '<input onclick="setmove(\''+movelist[k]+'\',\''+pokemonname+'\',\'1\',\''+channelset+'\');" class="btn btn-warning" type="button" value="Set 1st" id="setmove'+movelist[k]+'"/>';
										html += '<input onclick="setmove(\''+movelist[k]+'\',\''+pokemonname+'\',\'2\',\''+channelset+'\');" style="margin-left:5px" class="btn btn-warning" type="button" value="Set 2nd" id="setmove'+movelist[k]+'"/>';
										html += '<input onclick="setmove(\''+movelist[k]+'\',\''+pokemonname+'\',\'3\',\''+channelset+'\');" style="margin-left:5px" class="btn btn-warning" type="button" value="Set 3rd" id="setmove'+movelist[k]+'"/>';
										html += '<input onclick="setmove(\''+movelist[k]+'\',\''+pokemonname+'\',\'4\',\''+channelset+'\');" style="margin-left:5px" class="btn btn-warning" type="button" value="Set 4th" id="setmove'+movelist[k]+'"/>';
									html += '</form>';
									html += '<button class="'+lines[movelist[k]]['Type']+'" disabled>'+lines[movelist[k]]['Type']+'</button>';
									if(lines[movelist[k]]['Effect'] != ""){
										html += '<span class="label-text moves"> '+movelist[k]+' <i class="fa fa-info-circle" data-toggle="tooltip" title="" id="tootltip" data-original-title="'+lines[movelist[k]]['Effect']+'"></i> for <img src="./img/'+pokemonname+'.png" class="img-circle" width="50" height="50"> '+pokemonname+'</span>';
									}else{
										html += '<span class="label-text moves"> '+movelist[k]+' for <img src="./img/'+pokemonname+'.png" class="img-circle" width="50" height="50"> '+pokemonname+'</span>';
									};
								html += '</p>';
								document.getElementById("moveslist").innerHTML += html;
							}
						};
						$(function () {
							$('[data-toggle="tooltip"]').tooltip()
						})
				}).catch(err => console.error(err));
				break;
		};
	};
	ws.onclose = function() { 
		alert("Connection is closed..."); 
	};
} else {
	alert("Your browser is NOT supported!");
};
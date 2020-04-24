var clusterize;

window.onload = function() {
	if (typeof autocatcher !== 'undefined') {
		if (autocatcher == true) {
			document.getElementById("catcher").checked = true;
		};
	};
	if (typeof duplicate !== 'undefined') {
		if (duplicate == true) {
			document.getElementById("duplicate").checked = true;
		};
	};
	if (typeof aliases !== 'undefined') {
		if (aliases == true) {
			document.getElementById("aliases").checked = true;
		};
	};
	if (typeof farmer !== 'undefined') {
		farmtoken1 = farmer['token1'];
		farmtoken2 = farmer['token2'];
		farmchannelid = farmer['channelid'];
		farmstate = farmer['state'];
		if (farmtoken1 != 'none') {
			document.getElementById("farmtoken1").placeholder = farmtoken1;
		}
		if (farmtoken2 != 'none') {
			document.getElementById("farmtoken2").placeholder = farmtoken2;
		}
		if (farmchannelid != 'none') {
			document.getElementById("farmchannelid").placeholder = farmchannelid;
		}
		document.getElementById("farmactive").checked = farmstate;
	};
	if (token !== 'undefined') {
		document.getElementById("tokentext").placeholder = token;
	};
	if (autocatchdelay !== 'undefined') {
		document.getElementById("delayms").placeholder = autocatchdelay;
	};
	if (prefixes !== 'undefined') {
		document.getElementById("sbprefix").placeholder = prefixes["pokebot"];
		document.getElementById("pcprefix").placeholder = prefixes["pokecord"];
	};
	if (spamactive !== 'undefined') {
		if (spamactive == true) {
			document.getElementById("spammer").checked = spamactive;
		};
	};
	if (ClearAll(selected['name']) !== '') {
		document.getElementById("img").innerHTML = '<img src="./img/' + ClearAll(selected['name']) + '.png" class="img-circle" width="50" height="50">';
	};
	if (textchannel == 'undefined') {
		document.getElementById("alertbox").classList.add("alert");
		document.getElementById("alertbox").classList.add("alert-danger");
		document.getElementById("alertbox").innerHTML += '<p><label><span>You haven\'t registered a text channel where I can talk ! (p^register)</span></label></p>';
	} else {
		document.getElementById("alertbox").parentNode.removeChild(document.getElementById("alertbox"));
	};
	
	if(typeof listobj !== 'undefined'){
		LoadList(listobj);
	} else {
		LoadList(undefined);
	}
	
	if (servernames !== 'undefined') {
		serverarray = servernames.split(';');
		serveridarray = serverid.split(';');
		var html = '';
		whitelistobj = JSON.parse(whitelist);
		for (k = 0; k < serverarray.length; k++) {
			var checked = '';
			if (whitelistobj[serveridarray[k]]) {
				checked = 'checked';
			};
			if (serverarray[k] != '') {
				html = html + '<div class="form-check"><label><form><input ' + checked + ' onchange="refreshwhitelist(\'' + serveridarray[k] + '\');" type="checkbox" id="server' + serveridarray[k] + '"><span class="label-text"> ' + serverarray[k] + ' </span></form></label></div>'
			};
			if (k == serverarray.length - 1) {
				document.getElementById("listeserver").innerHTML = html;
			};
		};
	} else {
		var noservers = '<label><span style="color:red;" class="label-text">You don\'t have any servers :c</span></label>'
		document.getElementById("listeserver").innerHTML = noservers;
	};
};
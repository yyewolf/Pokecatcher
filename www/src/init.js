var clusterize;
var wlvalue;

window.onload = function() {
	wlvalue = ""
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
	if (typeof filter !== 'undefined') {
		if (filter == true) {
			document.getElementById("filter").checked = true;
		};
	};
	if (typeof customfilters !== 'undefined') {
		if (customfilters == true) {
			document.getElementById("customfilter").checked = true;
		};
	};
	if (token !== '') {
		document.getElementById("tokentext").placeholder = token;
	};
	if (queue !== '') {
		document.getElementById("alt").placeholder = queue;
	};
	if (autocatchdelay !== 'undefined') {
		document.getElementById("delayms").placeholder = autocatchdelay;
	};
	if (prefixes !== 'undefined') {
		document.getElementById("sbprefix").placeholder = prefixes["pokebot"];
		document.getElementById("pcprefix").placeholder = prefixes["pokecord"];
	};
	if (autolevelmax !== '') {
		document.getElementById("alm").placeholder = autolevelmax;
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
	LoadPokemons()
	FilterToHTML(filters)
};
var whitelistdisplay = 0

function LoadPokemons(){
	var html = [];
	names = Object.keys(pokewhitelist)
	for (k = 0; k < names.length; k++) {
		if(whitelistdisplay == 0 || (whitelistdisplay == 1 && pokewhitelist[names[k]])) {
			html.push(PokeHTML(names[k], pokewhitelist[names[k]]))
		}
	};
	clusterize = new Clusterize({
		rows: html,
		scrollId: 'WLScroll',
		contentId: 'WLContent'
	});
}

function PokeHTML(name, checked){
	t = ""
	if(checked){
		t = "checked"
	}
	
	return `
	<div id="wl_${name}" style="background-color: transparent">
		<label>
			<form style="float: center" id="pkmnwshitelist${name}">
				<input ${t} onchange="refreshpkmnwhitelist('${name.replace("'", "\\'")}',this.checked);" type="checkbox" id="wlpkmn${name}"/>
				<span class="label-text"><img src="./img/${name}.png" class="img-circle" width="50" height="50"> ${name}</span>
			</form>
		</label>
	</div>`
}

function Search(name) {
	if(wlvalue == name) {
		return
	}
	wlvalue = name
	var html = [];
	names = Object.keys(pokewhitelist)
	for (k = 0; k < names.length; k++) {
		if(names[k].includes(name)) {
			if(whitelistdisplay == 0 || (whitelistdisplay == 1 && pokewhitelist[names[k]])) {
				html.push(PokeHTML(names[k], pokewhitelist[names[k]]))
			}
		}
	};
	clusterize = new Clusterize({
		rows: html,
		scrollId: 'WLScroll',
		contentId: 'WLContent'
	});
}

function AllChecked(){
	names = Object.keys(pokewhitelist)
	for (k = 0; k < names.length; k++) {
		pokewhitelist[names[k]] = true
	};
	LoadPokemons()
	ws.send('{"action":"pkmnwhitelistallchecked"}');
};

function AllUnchecked(){
	names = Object.keys(pokewhitelist)
	for (k = 0; k < names.length; k++) {
		pokewhitelist[names[k]] = false
	};
	LoadPokemons()
	ws.send('{"action":"pkmnwhitelistallunchecked"}');
};

function LegendCheck(){
	names = Object.keys(pokewhitelist)
	for (k = 0; k < names.length; k++) {
		if(legendaries.includes(names[k])) {
			pokewhitelist[names[k]] = true
		}
	};
	LoadPokemons()
	ws.send('{"action":"pkmnwhitelistlegendchecked"}');
};

function download(data, filename, type) {
    var file = new Blob([data], {type: type});
    if (window.navigator.msSaveOrOpenBlob) // IE10+
        window.navigator.msSaveOrOpenBlob(file, filename);
    else { // Others
        var a = document.createElement("a"),
                url = URL.createObjectURL(file);
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        setTimeout(function() {
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);  
        }, 0); 
    }
}

function upload() {
    f = document.getElementById('imported')
	i = f.files.length - 1 
	file = f.files[i].text()
	file.then(function(val){
		pokewhitelist = JSON.parse(val)
		LoadPokemons()
		importpkmnwhitelist()
	})
}

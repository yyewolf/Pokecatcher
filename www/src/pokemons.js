function LoadList(list){
	if(typeof list !== 'undefined'){
		var html = [];
		for (k = 0; k < list['array']; k++) {
			if (typeof list[k + 1] != undefined) {
				try {
					var pokemonnumber = k + 1;
					var pokemonname = list[k + 1]['name']
					var pokemonlevel = list[k + 1]['level']
					var pokemoniv = list[k + 1]['iv']
				} catch (e) {};
				var ThisPkmn = new Pokemon(pokemonnumber, pokemonname, pokemonlevel, pokemoniv);
				html[k] = ThisPkmn.ListHTML
			};
		};
		clusterize = new Clusterize({
			rows: html,
			scrollId: 'PokeScroll',
			contentId: 'PokeContent'
		});
	}else{
		var nopokemons = '<label><span style="color:red;" class="label-text">You don\'t have any pokemons / Your list isn\'t loaded :c</span></label>'
		document.getElementById("PokeContent").innerHTML = nopokemons;
	};
}
function ReloadList(list){
	if(typeof list !== 'undefined'){
		var html = [];
		for (k = 0; k < list['array']; k++) {
			if (typeof list[k + 1] != undefined) {
				try {
					var pokemonnumber = k + 1;
					var pokemonname = list[k + 1]['name']
					var pokemonlevel = list[k + 1]['level']
					var pokemoniv = list[k + 1]['iv']
				} catch (e) {};
				var ThisPkmn = new Pokemon(pokemonnumber, pokemonname, pokemonlevel, pokemoniv);
				html[k] = ThisPkmn.ListHTML
			};
		};
		if(typeof clusterize !== 'undefined'){
			clusterize.update(html)
		}else {
			clusterize = new Clusterize({
				rows: html,
				scrollId: 'PokeScroll',
				contentId: 'PokeContent'
			});
		}
	}else{
		var nopokemons = '<label><span style="color:red;" class="label-text">You don\'t have any pokemons / Your list isn\'t loaded :c</span></label>'
		document.getElementById("PokeContent").innerHTML = nopokemons;
	};
}

function ClearAll(str) {
	str = str.replace(/\**/g, '');
	str = str.replace('**', '');
	str = str.replace(/ /g, '');
	str = str.replace(/é/g, 'e');
	str = str.replace(/è/g, 'e');
	str = str.replace('♂', '');
	str = str.replace('♀', '');
	return str
}

var legendaries = ['Moltres', 'Articuno', 'Zapdos', 'Mewtwo', 'Mew', 'Entei', 'Raikou', 'Suicune', 'Ho-Oh', 'Lugia', 'Celebi', 'Regice', 'Regirock', 'Registeel', 'Regigigas', 'Latios', 'Latias', 'Groudon', 'Kyogre', 'Rayquaza', 'Jirachi', 'Deoxys', 'Uxie', 'Mesprit', 'Azelf', 'Dialga', 'Palkia', 'Giratina', 'Arceus', 'Cresselia', 'Darkrai', 'Shaymin', 'Heatran', 'Manaphy', 'Phione', 'Cobalion', 'Terrakion', 'Virizion', 'Keldeo', 'Tornadus', 'Thundurus', 'Landorus', 'Zekrom', 'Reshiram', 'Kyurem', 'Victini', 'Meloetta', 'Genesect', 'Xerneas', 'Yveltal', 'Zygarde', 'Hoopa', 'Volcanion', 'Silvally', 'Tapu-Lele', 'Tapu-Koko', 'Tapu-Bulu', 'Tapu-Fini', 'Cosmog', 'Cosmoem', 'Solgaleo', 'Lunala', 'Necrozma', 'Magearna', 'Marshadow', 'Zeraora', 'Poipole', 'Naganadel', 'Guzzlord', 'Celesteela', 'Kartana', 'Xurkitree', 'Blacephalon', 'Buzzwole', 'Pheromosa', 'Nihilego', 'Stakataka']

class Pokemon {
  constructor(id, name, level, iv) {
    this._id = id;
    this._name = name;
    this._level = level;
    this._iv = iv;
	
	this._color = 'background-color: transparent';
	if (legendaries.includes(this._name)) {
		this._color = 'background-color: Gold';
	};
  }
 
  get ListHTML() {
    return `
	<div id="divpokemon${this._id}" style="${this._color}">
		<p style="display:inline-block;">
		<form style="float: left; padding: 5px;" id="formreleasepokemon${this._id}">
			<input onclick="dothis(\'release\',${this._id},\'${this._name}\');" class="btn btn-danger" type="button" value="Release" id="releasepokemon${this._id}"/>
			<input onclick="dothis(\'remove\',${this._id},\'${this._name}\');" style="margin-left:5px" class="btn btn-danger" type="button" value="Remove" id="remove${this._id}"/>
		</form>
		<form style="float: left; padding: 5px;" id="formselectpokemon${this._id}">
			<input onclick="dothis(\'select\',${this._id},\'${this._name}\');" class="btn btn-primary" type="button" value="Select" id="selectpokemon${this._id}"/>
		</form>
		<span class="label-text">Level ${this._level} <img src="./img/${this._name}.png" class="img-circle" width="50" height="50"> ${this._name} (${this._iv})</span>
		</p>
	</div>`
  }
}
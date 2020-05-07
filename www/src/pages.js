var page = 1;
		
function pages(nb){
	switch (nb) {
		case 1:
			document.getElementById("Pkmn").style["display"] = "contents";
			document.getElementById("SpammerP").style["display"] = "none";
			document.getElementById("PkmnWhitelist").style["display"] = "none";
			document.getElementById("CustomFilterP").style["display"] = "none";
			document.getElementById("OptionsP").style["display"] = "none";
			break;
		case 2:
			document.getElementById("Pkmn").style["display"] = "none";
			document.getElementById("SpammerP").style["display"] = "contents";
			document.getElementById("PkmnWhitelist").style["display"] = "none";
			document.getElementById("CustomFilterP").style["display"] = "none";
			document.getElementById("OptionsP").style["display"] = "none";
			break;
		case 3:
			document.getElementById("Pkmn").style["display"] = "none";
			document.getElementById("SpammerP").style["display"] = "none";
			document.getElementById("PkmnWhitelist").style["display"] = "contents";
			document.getElementById("CustomFilterP").style["display"] = "none";
			document.getElementById("OptionsP").style["display"] = "none";
			break;
		case 4:
			document.getElementById("Pkmn").style["display"] = "none";
			document.getElementById("SpammerP").style["display"] = "none";
			document.getElementById("PkmnWhitelist").style["display"] = "none";
			document.getElementById("CustomFilterP").style["display"] = "contents";
			document.getElementById("OptionsP").style["display"] = "none";
			break;
		case 5:
			document.getElementById("Pkmn").style["display"] = "none";
			document.getElementById("SpammerP").style["display"] = "none";
			document.getElementById("PkmnWhitelist").style["display"] = "none";
			document.getElementById("CustomFilterP").style["display"] = "none";
			document.getElementById("OptionsP").style["display"] = "contents";
			break;
	}
};
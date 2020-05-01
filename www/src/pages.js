var page = 1;
		
function pagenext(){
	page += 1;
	if(page == 4){
		page = 1;
		document.getElementById("firstPage").style["display"] = "contents";
		document.getElementById("secondPage").style["display"] = "none";
		document.getElementById("thirdPage").style["display"] = "none";
	}else if(page == 3){
		document.getElementById("firstPage").style["display"] = "none";
		document.getElementById("secondPage").style["display"] = "none";
		document.getElementById("thirdPage").style["display"] = "contents";
	}else if(page == 2){
		document.getElementById("firstPage").style["display"] = "none";
		document.getElementById("secondPage").style["display"] = "contents";
		document.getElementById("thirdPage").style["display"] = "none";
	}
};

function pageprevious(){
	page -= 1;
	if(page == 0){
		page = 3;
		document.getElementById("firstPage").style["display"] = "none";
		document.getElementById("secondPage").style["display"] = "none";
		document.getElementById("thirdPage").style["display"] = "contents";
	}else if(page == 2){
		document.getElementById("firstPage").style["display"] = "none";
		document.getElementById("secondPage").style["display"] = "contents";
		document.getElementById("thirdPage").style["display"] = "none";
	}else if(page == 1){
		document.getElementById("firstPage").style["display"] = "contents";
		document.getElementById("secondPage").style["display"] = "none";
		document.getElementById("thirdPage").style["display"] = "none";
	};
};
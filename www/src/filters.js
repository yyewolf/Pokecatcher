// Looks like this :
/* 
[
	{
		"Conditions":
			[
				{
					"Checking":1,
					"Operation":1,
					"Value":1
				},
				{
					"Checking":1,
					"Operation":1,
					"Value":1
				}
			]
	},
	{
		"Conditions":
			[
				{
					"Checking":1,
					"Operation":1,
					"Value":1
				},
				{
					"Checking":1,
					"Operation":1,
					"Value":1
				}
			]
	}
]
*/

function FilterToHTML(filters) {
	html = ""
	for(i = 0; i < filters.length; i++) {
		if (i != 0) {
			html += `<h2 style="font-size: 1rem;" class="label-text">AND</h2><br/>`
		}
		html += `
			<div id="FiltersArea${i}" class="col-lg-12" style="border:1px solid #00ff00;">
				<input type="button" onclick="AddCondition(${i});" class="btn btn-validate float-left" id="achecked" name="achecked" value="Add Condition" />
				<input type="button" onclick="RemoveFilter(${i});" class="btn btn-danger float-right" id="achecked" name="achecked" value="X" />
				<br /><br />`
		for(j = 0; j < filters[i].Conditions.length; j++) {
			current = filters[i].Conditions[j]
			if (j != 0) {
				html += `<h2 style="font-size: 1rem;" class="label-text">OR</h2><br/>`
			}
			html += `<div id="Filter${j}" class="col-lg-12" style="border:1px solid #cecece;">
						<input type="button" onclick="RemoveCondition(${i}, ${j});" class="btn btn-danger float-right" id="achecked" name="achecked" value="X" />
						<br />
						<p style="display:inline-block;">
							<form id="filter${i}${j}" style="padding: 8px;">
								<select onchange="UpdateFilters(${i}, ${j}, 'Checking', this.selectedIndex-1);" onload="$(this).find('option[value=${current.Checking}]').attr('selected','selected');" class="custom-select" style="width:180px; display:inline-block;" id="inputGroupSelect01">
									<option>Choose...</option>
									<option value="0">Level</option>
									<option value="1">HP (max 31)</option>
									<option value="2">Attack (max 31)</option>
									<option value="3">Defense (max 31)</option>
									<option value="4">SpAttack (max 31)</option>
									<option value="5">SpDef (max 31)</option>
									<option value="6">Speed (max 31)</option>
									<option value="7">IV</option>
									<option value="8">is Alolan</option>
									<option value="9">is Galarian</option>
									<option value="10">is Legendary</option>
									<option value="11">is Shiny</option>
								</select>
								<select onchange="UpdateFilters(${i}, ${j}, 'Operation', this.selectedIndex-1);" onload="this.selectedIndex=${current.Operation}+1" class="custom-select" style="width:120px; display:inline-block;" id="inputGroupSelect01">
									<option>Choose...</option>
									<option value="0">is under</option>
									<option value="1">is equal to</option>
									<option value="2">is above</option>
								</select>
								<input style="width:100px; display:inline-block;" type="number" class="form-control" id="idk" onkeyup="UpdateFilters(${i}, ${j}, 'Value', parseFloat(this.value));" value="${current.Value}">
							</form>
						</p>
						<br />
					</div>
					<br />` 
		}
		html += `</div><br />`
	}
	document.getElementById("FilterGoesHere").innerHTML = html;
	x = $(document.getElementById("FilterGoesHere")).find( "select" ).get()
	for(i = 0; i < x.length; i++) {
		try {
			x[i].onload()
		} catch(e) {
			continue
		}
	}
}

function AddFilter() {
	newObj = {
		"Conditions":[
			{
				"Checking":-1,
				"Operation":-1,
				"Value":50
			}
		]
	}
	filters.push(newObj)
	FilterToHTML(filters)
}

function AddCondition(i) {
	newObj = {
		"Checking":-1,
		"Operation":-1,
		"Value":50
	}
	filters[i].Conditions.push(newObj)
	FilterToHTML(filters)
}

function RemoveFilter(i) {
	filters.splice(i, 1);
	FilterToHTML(filters)
}

function RemoveCondition(i, j) {
	filters[i].Conditions.splice(j, 1);
	if(filters[i].Conditions.length == 0) {
		filters.splice(i, 1);
	}
	FilterToHTML(filters)
}

function UpdateFilters(i, j, c, v) {
	filters[i].Conditions[j][c] = v
	if(v >= 8) {
		x = document.getElementById("filter"+i+j)
		x.children[1].style= "width:120px; display:inline-block; visibility:hidden;"
		x.children[2].style= "width:100px; display:inline-block; visibility:hidden;"
	} else {
		x = document.getElementById("filter"+i+j)
		x.children[1].style= "width:120px; display:inline-block; visibility:content;"
		x.children[2].style= "width:100px; display:inline-block; visibility:content;"
	}
}
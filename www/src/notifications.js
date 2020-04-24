function notify(type, title = 0, icon = 0, message = 0, time = 5000, channelof = 0){
	if(type == 0){
		type = 'info';
	};
	
	template = '<div data-notify="container" class="col-xs-11 col-sm-3 alert alert-{0}" role="alert"><button type="button" aria-hidden="true" class="close" data-notify="dismiss">×</button>';
	
	if(icon != 0){
		template = template + '<img data-notify="icon" class="img-circle pull-left" style="padding-right:15px;"> ';
	};
		
	if(title != 0){
		template = template + '<span style="display: block;" data-notify="title">{1}</span> ';
	};
	
	if(message != 0 && channelof != 0){
		template = template + '<span style="margin-top: 30px;" data-notify="message">{2}</span>' +
							'<br />' +
							'<span style="margin-top: 30px;" data-notify="message">'+ channelof +
							'<div style="margin-top: 20px;" class="progress" data-notify="progressbar">' +
								'<div class="progress-bar progress-bar-{0}" role="progressbar" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100" style="width: 0%;"></div>' +
							'</div>' +
							'<a href="{3}" target="{4}" data-notify="url"></a>' +
						'</div>';
	}else if(message != 0){
		template = template + '<span style="margin-top: 30px;" data-notify="message">{2}</span>' +
							'<div style="margin-top: 20px;" class="progress" data-notify="progressbar">' +
								'<div class="progress-bar progress-bar-{0}" role="progressbar" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100" style="width: 0%;"></div>' +
							'</div>' +
							'<a href="{3}" target="{4}" data-notify="url"></a>' +
						'</div>';
	}else if(channelof != 0){
		template = template + '<span style="margin-top: 30px;" data-notify="message">' + channelof +
							'<div style="margin-top: 20px;" class="progress" data-notify="progressbar">' +
								'<div class="progress-bar progress-bar-{0}" role="progressbar" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100" style="width: 0%;"></div>' +
							'</div>' +
							'<a href="{3}" target="{4}" data-notify="url"></a>' +
						'</div>';
	}else{
		template = template + '<div style="margin-top: 20px;" class="progress" data-notify="progressbar">' +
								'<div class="progress-bar progress-bar-{0}" role="progressbar" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100" style="width: 0%;"></div>' +
							'</div>' +
							'<a href="{3}" target="{4}" data-notify="url"></a>' +
						'</div>';
	};
		
		
	$.notify({
		// options
		icon: icon,
		title: title,
		message: message
	},{
		// settings
		element: 'body',
		type: type,
		newest_on_top: true,
		showProgressbar: true,
		placement: {
			from: "top",
			align: "right"
		},
		delay: time,
		animate: {
			enter: 'animated fadeInRight',
			exit: 'animated fadeOutRight'
		},
		icon_type: 'image',
		template: template
	});
};


/*
	
	libUX.form HTML attributes:
		* data-form-target		Target (like form action)
		* data-form-method		Method	[POST, PUT, GET, DELETE]
		* data-form-trans-enc 			[JSON, JSONP, XML]
		* data-form-recv-enc 			[JSON, JSONP, XML]
*/

var libUX = {
	form: {
		getData: function(elem) {
			var data = {};
			var felems = elem.find("input, textarea, select");
			felems.each(function(ix, fe){
				var me = $(this);
				if (me.attr("name") != "") {	
					if (me.attr("type") == "checkbox") {
						data[me.attr("name")] = me.is(":checked");
					} else {
						data[me.attr("name")] = me.val();						
					}
				}
			});
			return data;
		},
		
		ajaxFormLoad: function(elem, url){
			if(typeof url !== "string" || url == ""){
				if(elem.attr("data-frmdest") != ""){
					var url = elem.attr("data-frmdest");
				}else{
					console.error("Missing data-frmdest destination URL for form!");
					return
				}
			}
			
			$.ajax({
				url: url, 
				//data: JSON.stringify(me.getData(elem)),
				method: "GET",
    			dataType: "json",
				success: function(opt){
					if(typeof opt.data == "object"){
						$.each(opt.data, function(k, v){
							var iel = elem.find('[name="' + k + '"]').first();
							if(iel.is('input') || iel.is('textarea')){
								iel.val(v);
							}
							/* else if(iel.is('select')){
								iel.val(v);
							}*/
						});
					}
				},
				error: function(opt){
					// !DUMMY
					console.log("Failed to load form data from service endpoint!");
				},
				async: true
			});
		},
		
		ajaxFormSubmit: function(elem, url, method, clb){
			var me = this;
			if(!method){
				var method = "POST";
			}
			if(typeof clb !== "function"){
				var clb = function(){
					// !DUMMY
					//libUX.message.show("Saved form!");
				};
			}
			$.ajax({
				url: url, 
				data: JSON.stringify(me.getData(elem)),
				method: "POST",
    			dataType: "json",
				success: function(opt){
					elem.find(".form-error").remove();
					if(typeof opt.meta.errors == "object" && Object.keys(opt.meta.errors).length >= 1) {
						$.each(opt.meta.errors, function(i, v){
							elem.find('[name=' + i + ']').after('<div class="form-error"></div>').next().text(v);
							// !DEBUG
							//console.log("Add elem: ", i, " -> ", v);
						});
					} else if (typeof clb === "function" && opt.meta.status == 200){
						clb(opt)
					}
				},
				error: function(opt){
					elem.find(".form-error").remove();
					elem.prepend('<div class="form-error"></div>').text(libUX.lang.form.err_submit);
				},
				async: true
			});
		}
	},
	message: {
		show: function(msg, opt){
			
		}
	},
	lang: {
		form: {
			err_submit: "Couldn't submit form! Error on webservice side."
		}
	}
};
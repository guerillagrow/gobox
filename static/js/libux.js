

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
					if (me.attr("type") == "checkbox" && me.is(":checked")) {
						data[me.attr("name")] = me.val();
					}else if (me.attr("type") == "radio" && me.is(":checked")) {
						data[me.attr("name")] = me.val();
					} else if (data[me.attr("name")] === undefined){
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
			
			// !TODO
			
			$.ajax({
				url: url, 
				//data: JSON.stringify(me.getData(elem)),
				method: "GET",
    			    dataType: "json",
				success: function(opt){
					if(typeof opt.data == "object"){
						$.each(opt.data, function(k, v){
							var iel = elem.find('[name="' + k + '"]').first();
							if(iel.is("input") && iel.attr("type") == "radio"){
								elem.find('input[name="' + k + '"][value="' + v + '"]').first().prop("checked", true);
							}else if(iel.is('input') || iel.is('textarea')){
								iel.val(v);
							}
							/* else if(iel.is('select')){
								iel.val(v);
							}*/
						});
					}
					if(typeof opt.meta == "object" && typeof opt.meta.__csrf__ == "string"){
						if (elem.find('input[name="__csrf__"][type="hidden"]').size() >= 1) {
							elem.find('input[name="__csrf__"][type="hidden"]').first().val(opt.meta.__csrf__);
						}else{
							var newEl = elem.append('<input type="hidden" name="__csrf__">');
							newEl.val(opt.meta.__csrf__)
						}
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
			var data = me.getData(elem);
			if (typeof data == "object" && typeof data.__csrf__ == "string") {
				url = libUX.url.addParam(url, "__csrf__", data.__csrf__);
				delete data["__csrf__"];
			}
			
			$.ajax({
				url: url, 
				data: JSON.stringify(data),
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
	
	url: {
		
	    params: function(url){
	
	        if(typeof url === 'string'){
	        
	            if(url.indexOf('?') > 0){
	                
	                var cutUrl = url.substr(url.indexOf('?') + 1);
	                var paramsRaw = cutUrl.split('&');
	                var paramList = {};
	                
	                
	                $.each(paramsRaw, function(k, v){
	                    
	                    var tmp = v.split('=');
	                    
	                    if(tmp[0] && tmp[1]){
	                        paramList[tmp[0]] = tmp[1];
	                    }else if(tmp[0]){
	                        paramList[tmp[0]] = '';
	                    }
	                    
	                });
	                
	                return paramList;
	                
	            }
	            
	        }
	        
	        return {};
	    },
		
	    addParam: function(url, param, val){
	        
			var me = this;
	        var paramsList = me.params(url);
	        var paramStr = '';
	        paramsList[param] = val;
	        
	        if(Object.keys(paramsList).length >= 1){
	            
	            var i = 0;
	            
	            $.each(paramsList, function(k, v){
	                paramStr = paramStr  + (i > 0 ? '&' : '') + encodeURIComponent(k) + '=' + encodeURIComponent(v);
	               i++;
	            });
	        
	        }
	        
	        var newUrl = url.substr(0, url.indexOf('?')) + (paramStr != '' ? '?' + paramStr : '');
			return newUrl;
	        
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
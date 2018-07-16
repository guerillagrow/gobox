

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
			
			// !TODO: fix param handling if no "?" inside url
			// !DEBUG
	
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
	        
			// !TODO: fix param handling if no "?" inside url
			// !DEBUG
			
			var me = this;
	        var paramsList = me.params(url);
	        var paramStr = '';
	        paramsList[param] = val;
			var index = url.indexOf('?');
			
			if(index == -1) {
				index = 0;
			}
	        
	        if(Object.keys(paramsList).length >= 1){	            
	            var i = 0;
	            
	            $.each(paramsList, function(k, v){
	                paramStr = paramStr  + (i > 0 ? '&' : '') + encodeURIComponent(k) + '=' + encodeURIComponent(v);
	               i++;
	            });	        
	        }
	        
	        var newUrl = url.substr(0, index) + (paramStr != '' ? '?' + paramStr : '');
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


String.prototype.replaceAll = function (find, replace) {
    var str = this;
    return str.replace(new RegExp(find, 'g'), replace);
};
Array.prototype.count = function () {
	for (var i = 0; i < this.length; i++) {
	   return this.length;
	}
};

Array.prototype.size = function () {
	for (var i = 0; i < this.length; i++) {
	   return this.length;
	}
};


/**********************************************
 * Date Format 1.2.3
 * Function: dateFormat / .format()
 * ------------------------------------------
 * 
 * (c) 2007-2009 Steven Levithan <stevenlevithan.com>
 * Includes enhancements by Scott Trenda <scott.trenda.net>
 * and Kris Kowal <cixar.com/~kris.kowal/>
 **********************************************/
var dateFormat = function () {
	var	token = /d{1,4}|m{1,4}|yy(?:yy)?|([HhMsTt])\1?|[LloSZ]|"[^"]*"|'[^']*'/g,
		timezone = /\b(?:[PMCEA][SDP]T|(?:Pacific|Mountain|Central|Eastern|Atlantic) (?:Standard|Daylight|Prevailing) Time|(?:GMT|UTC)(?:[-+]\d{4})?)\b/g,
		timezoneClip = /[^-+\dA-Z]/g,
		pad = function (val, len) {
			val = String(val);
			len = len || 2;
			while (val.length < len) val = "0" + val;
			return val;
		};

	// Regexes and supporting functions are cached through closure
	return function (date, mask, utc) {
		var dF = dateFormat;

		// You can't provide utc if you skip other args (use the "UTC:" mask prefix)
		if (arguments.length == 1 && Object.prototype.toString.call(date) == "[object String]" && !/\d/.test(date)) {
			mask = date;
			date = undefined;
		}

		// Passing date through Date applies Date.parse, if necessary
		date = date ? new Date(date) : new Date;
		if (isNaN(date)) throw SyntaxError("invalid date");

		mask = String(dF.masks[mask] || mask || dF.masks["default"]);

		// Allow setting the utc argument via the mask
		if (mask.slice(0, 4) == "UTC:") {
			mask = mask.slice(4);
			utc = true;
		}

		var	_ = utc ? "getUTC" : "get",
			d = date[_ + "Date"](),
			D = date[_ + "Day"](),
			m = date[_ + "Month"](),
			y = date[_ + "FullYear"](),
			H = date[_ + "Hours"](),
			M = date[_ + "Minutes"](),
			s = date[_ + "Seconds"](),
			L = date[_ + "Milliseconds"](),
			o = utc ? 0 : date.getTimezoneOffset(),
			flags = {
				d:    d,
				dd:   pad(d),
				ddd:  dF.i18n.dayNames[D],
				dddd: dF.i18n.dayNames[D + 7],
				m:    m + 1,
				mm:   pad(m + 1),
				mmm:  dF.i18n.monthNames[m],
				mmmm: dF.i18n.monthNames[m + 12],
				yy:   String(y).slice(2),
				yyyy: y,
				h:    H % 12 || 12,
				hh:   pad(H % 12 || 12),
				H:    H,
				HH:   pad(H),
				M:    M,
				MM:   pad(M),
				s:    s,
				ss:   pad(s),
				l:    pad(L, 3),
				L:    pad(L > 99 ? Math.round(L / 10) : L),
				t:    H < 12 ? "a"  : "p",
				tt:   H < 12 ? "am" : "pm",
				T:    H < 12 ? "A"  : "P",
				TT:   H < 12 ? "AM" : "PM",
				Z:    utc ? "UTC" : (String(date).match(timezone) || [""]).pop().replace(timezoneClip, ""),
				o:    (o > 0 ? "-" : "+") + pad(Math.floor(Math.abs(o) / 60) * 100 + Math.abs(o) % 60, 4),
				S:    ["th", "st", "nd", "rd"][d % 10 > 3 ? 0 : (d % 100 - d % 10 != 10) * d % 10]
			};

		return mask.replace(token, function ($0) {
			return $0 in flags ? flags[$0] : $0.slice(1, $0.length - 1);
		});
	};
}();
var Base64 = {

	// private property
	_keyStr : "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",

	// public method for encoding
	encode : function (input) {
		var output = "";
		var chr1, chr2, chr3, enc1, enc2, enc3, enc4;
		var i = 0;

		input = Base64._utf8_encode(input);

		while (i < input.length) {

			chr1 = input.charCodeAt(i++);
			chr2 = input.charCodeAt(i++);
			chr3 = input.charCodeAt(i++);

			enc1 = chr1 >> 2;
			enc2 = ((chr1 & 3) << 4) | (chr2 >> 4);
			enc3 = ((chr2 & 15) << 2) | (chr3 >> 6);
			enc4 = chr3 & 63;

			if (isNaN(chr2)) {
				enc3 = enc4 = 64;
			} else if (isNaN(chr3)) {
				enc4 = 64;
			}

			output = output +
			this._keyStr.charAt(enc1) + this._keyStr.charAt(enc2) +
			this._keyStr.charAt(enc3) + this._keyStr.charAt(enc4);

		}

		return output;
	},

	// public method for decoding
	decode : function (input) {
		var output = "";
		var chr1, chr2, chr3;
		var enc1, enc2, enc3, enc4;
		var i = 0;

		input = input.replace(/[^A-Za-z0-9\+\/\=]/g, "");

		while (i < input.length) {

			enc1 = this._keyStr.indexOf(input.charAt(i++));
			enc2 = this._keyStr.indexOf(input.charAt(i++));
			enc3 = this._keyStr.indexOf(input.charAt(i++));
			enc4 = this._keyStr.indexOf(input.charAt(i++));

			chr1 = (enc1 << 2) | (enc2 >> 4);
			chr2 = ((enc2 & 15) << 4) | (enc3 >> 2);
			chr3 = ((enc3 & 3) << 6) | enc4;

			output = output + String.fromCharCode(chr1);

			if (enc3 != 64) {
				output = output + String.fromCharCode(chr2);
			}
			if (enc4 != 64) {
				output = output + String.fromCharCode(chr3);
			}

		}

		output = Base64._utf8_decode(output);

		return output;

	},

	// private method for UTF-8 encoding
	_utf8_encode : function (string) {
		string = string.replace(/\r\n/g,"\n");
		var utftext = "";

		for (var n = 0; n < string.length; n++) {

			var c = string.charCodeAt(n);

			if (c < 128) {
				utftext += String.fromCharCode(c);
			}
			else if((c > 127) && (c < 2048)) {
				utftext += String.fromCharCode((c >> 6) | 192);
				utftext += String.fromCharCode((c & 63) | 128);
			}
			else {
				utftext += String.fromCharCode((c >> 12) | 224);
				utftext += String.fromCharCode(((c >> 6) & 63) | 128);
				utftext += String.fromCharCode((c & 63) | 128);
			}

		}

		return utftext;
	},

	// private method for UTF-8 decoding
	_utf8_decode : function (utftext) {
		var string = "";
		var i = 0;
		var c = c1 = c2 = 0;

		while ( i < utftext.length ) {

			c = utftext.charCodeAt(i);

			if (c < 128) {
				string += String.fromCharCode(c);
				i++;
			}
			else if((c > 191) && (c < 224)) {
				c2 = utftext.charCodeAt(i+1);
				string += String.fromCharCode(((c & 31) << 6) | (c2 & 63));
				i += 2;
			}
			else {
				c2 = utftext.charCodeAt(i+1);
				c3 = utftext.charCodeAt(i+2);
				string += String.fromCharCode(((c & 15) << 12) | ((c2 & 63) << 6) | (c3 & 63));
				i += 3;
			}

		}

		return string;
	}

};
String.prototype.replaceAll = function (find, replace) {
    var str = this;
    return str.replace(new RegExp(find, 'g'), replace);
};
Array.prototype.count = function () {
	for (var i = 0; i < this.length; i++) {
	   return this.length;
	}
}

Array.prototype.size = function () {
	for (var i = 0; i < this.length; i++) {
	   return this.length;
	}
}
/*
Object.prototype.size = function(obj) {
    var size = 0, key;
    for (key in obj) {
        if (obj.hasOwnProperty(key)) size++;
    }
    return size;
};*/

function objCount(obj){
	var count = 0;
	for (var k in obj) {
	  // if the object has this property and it isn't a property
	  // further up the prototype chain
	  if (obj.hasOwnProperty(k)) count++;
	}
	return count;
}
// Some common format strings
dateFormat.masks = {
	"default":      "ddd mmm dd yyyy HH:MM:ss",
	shortDate:      "m/d/yy",
	mediumDate:     "mmm d, yyyy",
	longDate:       "mmmm d, yyyy",
	fullDate:       "dddd, mmmm d, yyyy",
	shortTime:      "h:MM TT",
	mediumTime:     "h:MM:ss TT",
	longTime:       "h:MM:ss TT Z",
	isoDate:        "yyyy-mm-dd",
	isoTime:        "HH:MM:ss",
	isoDateTime:    "yyyy-mm-dd'T'HH:MM:ss",
	isoUtcDateTime: "UTC:yyyy-mm-dd'T'HH:MM:ss'Z'"
};

// Internationalization strings
dateFormat.i18n = {
	dayNames: [
		"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat",
		"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"
	],
	monthNames: [
		"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
		"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"
	]
};

// PROTOTYPES
Date.prototype.format = function (mask, utc) {
	return dateFormat(this, mask, utc);
};
Date.prototype.daysInMonth = function(month, year) {
    var month = month || this.getMonth() + 1;
    var year = year || this.getFullYear();
    var dd = new Date(year, month, 0);
    return dd.getDate();
}
Date.prototype.toJson = function () {
    return "\/Date(" + +new Date(Date.UTC(
        this.getFullYear(),
        this.getMonth(),
        this.getDate(),
        this.getHours(),
        this.getMinutes()
    )) + ")\/";
};


function isObject(ivar){

    if(ivar !== null && typeof ivar === 'object'){
        return true;
    }else{
        return false;
    }

}

function isArray(ivar){

    if(ivar !== null && typeof ivar === 'array'){
        return true;
    }else{
        return false;
    }

}


function isNumber(ivar){

    if(ivar !== null && typeof ivar === 'number'){
        return true;
    }else{
        return false;
    }

}

function isBool(ivar){

    if(typeof ivar === 'boolean'){
        return true;
    }else{
        return false;
    }

}

function isString(ivar){

    if(typeof ivar === 'string'){
        return true;
    }else{
        return false;
    }

}

function isNull(ivar){

    if(typeof ivar === null){
        return true;
    }else{
        return false;
    }

}

function isUndefined(ivar){

    if(typeof ivar === 'undefined'){
        return true;
    }else{
        return false;
    }

}

function isFunction(ivar){

    if(typeof ivar === 'function'){
        return true;
    }else{
        return false;
    }

}
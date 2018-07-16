
<!DOCTYPE html>
<html>
  <head>
    <title>GoBox</title>
<link rel="shortcut icon" type="image/png" href="/static/img/cannabis-logo-sm_x16.png"/>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <!-- jQuery UI -->
    <!-- link href="https://code.jquery.com/ui/1.10.3/themes/redmond/jquery-ui.css" rel="stylesheet" media="screen" -->

    <!-- Bootstrap -->
    <link href="/static/bootstrap/css/bootstrap.min.css" rel="stylesheet">
    <!-- styles -->
    <link href="/static/css/styles.css" rel="stylesheet">

    <link href="/static/css/stats.css" rel="stylesheet">

    <!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
      <script src="https://oss.maxcdn.com/libs/respond.js/1.3.0/respond.min.js"></script>
    <![endif]-->
	
	
	<style type="text/css">
		.panel-options i.glyphicon {
		    font-size: 24px;  
		}
		form .form-error {
			margin: 5px;
			display:block;
			padding: 10px;
			border: 2px solid red;
		}
	</style>
	
	
  </head>

  <body>
  	<div class="header" style="min-height:70px;">
	     <div class="container">
	        <div class="row">
	           <div class="col-md-12">
	              <!-- Logo -->
	              <div class="logo">
	                 <h1 style="width:125px;display:inline-block;"><a href="/" style="display:inline;">GoBo<img src="/static/img/cannabis-logo-sm.png" style="width:40px;display:inline;margin-top:-10px;margin-left:-5px;"></a></h1>
					
	              </div>
	           </div>
	           <div class="col-md-5" style="display:none;">
	              <div class="row">
	                <div class="col-lg-12">
					
	                  <!--div class="input-group form">
	                       <input type="text" class="form-control" placeholder="Search...">
	                       <span class="input-group-btn">
	                         <button class="btn btn-primary" type="button">Search</button>
	                       </span>
	                  </div -->
	                </div>
	              </div>
	           </div>
	           <div class="col-md-2" style="display:none;">
	              <div class="navbar navbar-inverse" role="banner">
	                  <nav class="collapse navbar-collapse bs-navbar-collapse navbar-right" role="navigation">
	                    <ul class="nav navbar-nav">
	                      <li class="dropdown">
	                        <a href="#" class="dropdown-toggle" data-toggle="dropdown">My Account <b class="caret"></b></a>
	                        <ul class="dropdown-menu animated fadeInUp">
	                          <li><a href="profile.html">Profile</a></li>
	                          <li><a href="login.html">Logout</a></li>
	                        </ul>
	                      </li>
	                    </ul>
	                  </nav>
	              </div>
	           </div>
	        </div>
	     </div>
	</div>

    <div class="page-content">
		
	
		{{if  or (.sensor_t1) (.sensor_t2)}}
    		<div class="row">
		
		  <div class="col-md-12">

  			<div class="row">
				{{if .sensor_t1}}
  				<div class="col-md-6">
  					<div class="content-box-large" id="sensor-t1-chart">
		  				<div class="panel-heading">
							<div class="panel-title">Sensor T1</div>
							
							<div class="panel-options">
								<a href="#" data-rel="collapse" class="x-refresh"><i class="glyphicon glyphicon-refresh"></i></a>

							</div>
						</div>
		  				<div class="panel-body">
							<div>
								<b>Temperature:</b> <span class="cur-temp"></span> °C<br>
								<b>Humidity:</b> <span class="cur-hum"></span> % rH<br>
								<button class="btn btn-default x-tl-day x-tgl" data-value="day">Last Day</button> <button class="btn btn-default x-tl-hour x-tgl active"  data-value="hour">Last Hours</button>
							</div><br>
		  					<div class="tchart" style="width:100%;height:300px"></div>
		  				</div>
		  			</div>
  				</div>
  				{{end}}
				{{if .sensor_t2}}
				<div class="col-md-6">
  					<div class="content-box-large"  id="sensor-t2-chart">
		  				<div class="panel-heading">
							<div class="panel-title">Sensor T2</div>
							
							<div class="panel-options">
								<a href="#" data-rel="collapse" class="x-refresh"><i class="glyphicon glyphicon-refresh"></i></a>
							</div>
						</div>
		  				<div class="panel-body">
							<div>
								<b>Temperature:</b> <span class="cur-temp"></span> °C<br>
								<b>Humidity:</b> <span class="cur-hum"></span> % rH<br>
								<button class="btn btn-default x-tl-day x-tgl" data-value="day">Last Day</button> <button class="btn btn-default x-tl-hour x-tgl active"  data-value="hour">Last Hours</button>
							</div><br>
		  					<div class="tchart" style="width:100%;height:300px"></div>
		  				</div>
		  			</div>
  				</div>
  				{{end}}
			</div>

		  </div>
		</div>
		{{end}}
		
		{{if  or (.relay_l1) (.relay_l2)}}
    		<div class="row">
		
		  <div class="col-md-12">

  			<div class="row">
				{{if .relay_l1}}
  				<div class="{{if .relay_l2}}col-md-6{{else}}col-md-12{{end}} col-sm-12">
  					<div class="content-box-large">
		  				<div class="panel-heading">
							<div class="panel-title">Relay L1 Config</div>
							
							<!--div class="panel-options" style="display:none;">
								<a href="#" data-rel="collapse"><i class="glyphicon glyphicon-refresh"></i></a>
								<a href="#" data-rel="reload"><i class="glyphicon glyphicon-cog"></i></a>
							</div-->
						</div>
		  				<div class="panel-body">
							<form class="form-horizontal" role="form" id="svc-relay-l1-form" data-frmdest="/svc/relay?target=l1" data-frmdata="">
							  <div class="form-group">
							    <label for="ton" class="col-sm-2 control-label">Time On</label>
							    <div class="col-sm-10">
							      <input type="text" class="form-control" name="ton" placeholder="On-Time">
							    </div>
							  </div>
							  <div class="form-group">
							    <label for="toff" class="col-sm-2 control-label">Time Off</label>
							    <div class="col-sm-10">
							      <input type="text" class="form-control" name="toff" placeholder="Off-Time">
							    </div>
							  </div>
							  <div class="form-group">
							    <label for="cond" class="col-sm-2 control-label">Condition</label>
							    <div class="col-sm-10">
							      <input type="text" class="form-control" name="cond" placeholder="Relay switch condition">
							    </div>
							  </div>
							  <div class="form-group" style="">
							    <label for="cond" class="col-sm-2 control-label">Force</label>
							    <div class="col-sm-10">
							      <div class="checkbox">									
							        <label>
							          <!--input type="checkbox" name="status"> Relay Status -->
								        <input type="radio" value="-1" name="force"> Off
							        </label>
									
							        <label>
							          <!--input type="checkbox" name="status"> Relay Status -->
								        <input type="radio" value="1" name="force"> On
							        </label>
									
							        <label>
							          <!--input type="checkbox" name="status"> Relay Status -->
								        <input type="radio" value="0" name="force"> None
							        </label>
									
									
							      </div>
							    </div>
								
								<!--div class="btn-group btn-group-toggle" data-toggle="buttons">
								  <label class="btn btn-secondary active">
								    <input type="radio" name="force" autocomplete="off" value="1"> On
								  </label>
								  <label class="btn btn-secondary">
								    <input type="radio" name="force" autocomplete="off" value="-1"> Off
								  </label>
								  <label class="btn btn-secondary">
								    <input type="radio" name="force" autocomplete="off" value="0"> None
								  </label>
								</div-->
							  
							
							</div>
							  <div class="form-group">
							    <div class="col-sm-offset-2 col-sm-10">
							      <button type="submit" class="btn btn-primary">Save</button>
								  <input type="hidden" name="__csrf__" value="">
							    </div>
							  </div>
							</form>
		  				</div>
		  			</div>
  				</div>
				{{end}}
				{{if .relay_l2}}
  				<div class="{{if .relay_l1}}col-md-6{{else}}col-md-12{{end}} col-sm-12">
  					<div class="content-box-large">
		  				<div class="panel-heading">
							<div class="panel-title">Relay L2 Config</div>
							
							<!--div class="panel-options" style="display:none;">
								<a href="#" data-rel="collapse"><i class="glyphicon glyphicon-refresh"></i></a>
								<a href="#" data-rel="reload"><i class="glyphicon glyphicon-cog"></i></a>
							</div-->
						</div>
		  				<div class="panel-body">
							<form class="form-horizontal" role="form" id="svc-relay-l2-form" data-frmdest="/svc/relay?target=l2" data-frmdata="">
							  <div class="form-group">
							    <label for="ton" class="col-sm-2 control-label">Time On</label>
							    <div class="col-sm-10">
							      <input type="text" class="form-control" name="ton" placeholder="On-Time">
							    </div>
							  </div>
							  <div class="form-group">
							    <label for="toff" class="col-sm-2 control-label">Time Off</label>
							    <div class="col-sm-10">
							      <input type="text" class="form-control" name="toff" placeholder="Off-Time">
							    </div>
							  </div>
							  <div class="form-group">
							    <label for="cond" class="col-sm-2 control-label">Condition</label>
							    <div class="col-sm-10">
							      <input type="text" class="form-control" name="cond" placeholder="Relay switch condition">
							    </div>
							  </div>
							  <div class="form-group" style="">
							
							    <label for="cond" class="col-sm-2 control-label">Force</label>
							    <div class="col-sm-10">
							      <div class="checkbox">									
							        <label>
							          <!--input type="checkbox" name="status"> Relay Status -->
								        <input type="radio" value="-1" name="force"> Off
							        </label>
									
							        <label>
							          <!--input type="checkbox" name="status"> Relay Status -->
								        <input type="radio" value="1" name="force"> On
							        </label>
									
							        <label>
							          <!--input type="checkbox" name="status"> Relay Status -->
								        <input type="radio" value="0" name="force"> None
							        </label>
									
									
							      </div>
							    </div>
								
							  </div>
							  <div class="form-group">
							    <div class="col-sm-offset-2 col-sm-10">
							      <button type="submit" class="btn btn-primary">Save</button>
								  <input type="hidden" name="__csrf__" value="">
							    </div>
							  </div>
							</form>
		  				</div>
		  			</div>
  				</div>				
				{{end}}
			</div>
			</div>
		</div>
		{{end}}
    </div>

	<div class="row">
	
		<div class="col-md-12">
  				<div class="col-md-12">
  					<div class="content-box-large">
		  				<div class="panel-heading">
							<div class="panel-title">Account settings</div>
							
							<!--div class="panel-options" style="display:none;">
								<a href="#" data-rel="collapse"><i class="glyphicon glyphicon-refresh"></i></a>
								<a href="#" data-rel="reload"><i class="glyphicon glyphicon-cog"></i></a>
							</div-->
						</div>
		  				<div class="panel-body">
							<form class="form-horizontal" role="form" id="svc-user-form" data-frmdest="/svc/user" data-frmdata="">
							  <div class="form-group">
							    <label for="name" class="col-sm-2 control-label">Name</label>
							    <div class="col-sm-10">
							      <input type="text" class="form-control" name="name" placeholder="Username">
							    </div>
							  </div>
							  <div class="form-group">
							    <label for="email" class="col-sm-2 control-label">Email</label>
							    <div class="col-sm-10">
							      <input type="text" class="form-control" name="email" placeholder="Email">
							    </div>
							  </div>
							  <div class="form-group">
							    <label for="current_password" class="col-sm-2 control-label">Current Password</label>
							    <div class="col-sm-10">
							      <input type="password" class="form-control" name="current_password" placeholder="Current password">
							    </div>
							  </div>
							  <div class="form-group">
							    <label for="password" class="col-sm-2 control-label">New Password</label>
							    <div class="col-sm-10">
							      <input type="password" class="form-control" name="password" placeholder="Password">
							    </div>
							  </div>
							  <div class="form-group">
							    <div class="col-sm-offset-2 col-sm-10">
							      <button type="submit" class="btn btn-primary">Save</button>
								  <input type="hidden" name="__csrf__" value="">
							    </div>
							  </div>
							</form>
		  				</div>
		  			</div>
  				</div>
		</div>
	</div>
	
	
    <footer>
         <div class="container">
         
            <div class="copy text-center">
               Powered by <a target="_blank" href='https://github.com/guerillagrow/gobox'>GoBox</a> - the open source GrowBox automation toolkit!
            </div>
            
         </div>
      </footer>

    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <script src="/static/js/jquery.js"></script>
    <!-- jQuery UI -->
    <!--script src="https://code.jquery.com/ui/1.10.3/jquery-ui.js"></script-->
  
  <!-- Include all compiled plugins (below), or include individual files as needed -->
    <script src="/static/bootstrap/js/bootstrap.min.js"></script>

    <link rel="stylesheet" href="/static/vendors/morris/morris.css">


    <!--script src="/static/vendors/jquery.knob.js"></script-->
    <!--script src="/static/vendors/raphael-min.js"></script-->
    <script src="/static/vendors/morris/morris.min.js"></script>

    <script src="/static/vendors/flot/jquery.flot.js"></script>
    <script src="/static/vendors/flot/jquery.flot.categories.js"></script>
    <script src="/static/vendors/flot/jquery.flot.pie.js"></script>
    <script src="/static/vendors/flot/jquery.flot.time.js"></script>
    <script src="/static/vendors/flot/jquery.flot.stack.js"></script>
    <script src="/static/vendors/flot/jquery.flot.resize.js"></script>

	
	<link rel="stylesheet" type="text/css" href="/static/vendors/jquery.jgrowl.min.css" />
	<script src="/static/vendors/jquery.jgrowl.min.js"></script>
	
    <script src="/static/js/libux.js"></script>
    <script src="/static/js/custom.js"></script>
    <!--script src="/static/js/stats.js"></script-->
	<script>
		
		function doPlot(id, position, lt, incr, datax, datay) {
			if (incr < 2) {
				return
			}
			
			if ($(window).width() < 650) {
				
				if (lt == "day") {
					var cp = [5, "hour"];	
					var pc = {"show": false};
				}else{
					var cp = [60, "minute"];	
					var pc = {"show": false};
				}
			}else{
				if (lt == "day") {
					var cp = [4, "hour"];	
					var pc = {"show": false};
				}else{
					var cp = [15, "minute"];	
					var pc = {"show": false};
				}
				
			}
			
		
			$(id).find(".cur-temp").html(datax[0][1]);
			$(id).find(".cur-hum").html(datay[0][1]);
	    		$.plot(id + " .tchart", 
				[
			        { data: datax, label: "Temperature" },
			        { data: datay, label: "Humidity", yaxis: 2 }
			    ], {
					 yaxis: {
						min: 0
			        },
			        xaxis: { mode: "time",
							timeformat: "%a %H:%M:%S",
          					timezone: "browser", // localtime -> this converts UTC to  local date
							minTickSize: cp,
							//tickLength: 5,
							dayNames: ["So", "Mo", "Di", "Mi", "Do", "Fr", "Sa"],
							monthNames: ["Jan", "Feb", "Mär", "Apr", "Mai", "Jun", "Jul", "Aug", "sep", "okt", "nov", "dec"],
							//timeformat: "%Y-%m-%d %H:%M:%S",
			                //min: (new Date("2000/01/01")).getTime(),
			                //max: (new Date("2000/01/02")).getTime(),
							//alignTicksWithAxis: position == "right" ? 1 : null,
							//position: position
						alignTicksWithAxis: position == "right" ? 1 : null,
						position: position
					},
			       /* xaxes: [ { mode: "time" } ],
					
					
			        yaxes: [ { min: 0 }, {
			            // align if we are to the right
			            alignTicksWithAxis: position == "right" ? 1 : null,
			            position: position,
			            //tickFormatter: euroFormatter
			        } ],*/
			        legend: { position: "sw" },
					
     				grid:      { hoverable: true, clickable: false },
			        points: pc,
     				lines:  { show: true, lineWidth: 1 }, 
			        clickable:true,
					hoverable: true
		    });
		}
		
		var limitRq = 5000;
		
		function getStatsDataT2(tl, clb){
			var internalRes = {
				t1: {
					incr:0,
					temp: {},
					hum: {}
				},
				t2: {
					incr:0,
					temp: {},
					hum: {}
				}
			};
			$.ajax({url: "/svc/sensors/temperature?sensor=T2&g=1&limit="+limitRq, 
				data: {
					"tl": tl
				},
				success: function(opt){
					console.log(opt);
				
					internalRes["t2"]["incr"]++;
					internalRes["t2"]["temp"] = opt.data;
					clb(internalRes["t2"])
					//doPlot("#sensor-t2-chart","right", internalRes["t2"]["incr"], internalRes["t2"]["temp"], internalRes["t2"]["hum"]);
				},
				async: true
			});
			$.ajax({url: "/svc/sensors/humidity?sensor=T2&g=1&limit="+limitRq, 
				data: {
					"tl": tl
				},
				success: function(opt){
					console.log(opt);
					internalRes["t2"]["incr"]++;
					internalRes["t2"]["hum"] = opt.data;
					clb(internalRes["t2"])
					//doPlot("#sensor-t2-chart","right", internalRes["t2"]["incr"], internalRes["t2"]["temp"], internalRes["t2"]["hum"]);
				},
				async: true
			});
			return internalRes["t2"];
		}
		
		function getStatsDataT1(tl, clb){
			var internalRes = {
				t1: {
					incr:0,
					temp: {},
					hum: {}
				},
				t2: {
					incr:0,
					temp: {},
					hum: {}
				}
			};
			$.ajax({url: "/svc/sensors/temperature?sensor=T1&g=1&limit="+limitRq, 
				data: {
					"tl": tl
				},
				success: function(opt){
					console.log(opt);
				
				
					internalRes["t1"]["incr"]++;
					internalRes["t1"]["temp"] = opt.data;
					clb(internalRes["t1"])
					//doPlot("#sensor-t1-chart","right", internalRes["t1"]["incr"], internalRes["t1"]["temp"], internalRes["t1"]["hum"]);
				},
				async: true
			});
			$.ajax({url: "/svc/sensors/humidity?sensor=T1&g=1&limit="+limitRq, 
				data: {
					"tl": tl
				},
				success: function(opt){
					console.log(opt);
					internalRes["t1"]["incr"]++;
					internalRes["t1"]["hum"] = opt.data;
					clb(internalRes["t1"])
					//doPlot("#sensor-t1-chart","right", internalRes["t1"]["incr"], internalRes["t1"]["temp"], internalRes["t1"]["hum"]);
				},
				async: true
			});
			return internalRes["t1"];
			
			
			
			
		}
		
		function __init__(){
			{{if .sensor_t1}}
			getStatsDataT1("hour", function(d){
				doPlot("#sensor-t1-chart","right", "hour", d["incr"], d["temp"], d["hum"]);
			});
			{{end}}
			{{if .sensor_t2}}
			getStatsDataT2("hour", function(d){
				doPlot("#sensor-t2-chart","right","hour", d["incr"], d["temp"], d["hum"]);
			});
			{{end}}
			{{if .relay_l1}}
			libUX.form.ajaxFormLoad($("#svc-relay-l1-form"));
			{{end}}
			{{if .relay_l2}}
			libUX.form.ajaxFormLoad($("#svc-relay-l2-form"));
			{{end}}
			libUX.form.ajaxFormLoad($("#svc-user-form"));
		}
		
		$(document).ready(function(){
			__init__();
			$("#sensor-t1-chart .x-refresh").on("click", function(e){
				e.preventDefault();
				getStatsDataT1($("#sensor-t1-chart .x-tgl.active").attr("data-value"), function(d){
					if(d["incr"] == 2){
						doPlot("#sensor-t1-chart","right", $("#sensor-t1-chart .x-tgl.active").attr("data-value"), d["incr"], d["temp"], d["hum"]);
						$.jGrowl("Got sensor data", { 
							life: 1000, 
							closerTemplate: "<div>[ close all ]</div>",
							closeTemplate: "×" 
						});
					}
				});
			});
			$("#sensor-t2-chart .x-refresh").on("click", function(e){
				e.preventDefault();
				getStatsDataT2($("#sensor-t2-chart .x-tgl.active").attr("data-value"), function(d){
					if(d["incr"] == 2){
						doPlot("#sensor-t2-chart","right",$("#sensor-t2-chart .x-tgl.active").attr("data-value"), d["incr"], d["temp"], d["hum"]);
						$.jGrowl("Got sensor data", { 
							life: 1000, 
							closerTemplate: "<div>[ close all ]</div>",
							closeTemplate: "×" 
						});
					}
				});
			});
			
			$("#sensor-t1-chart .x-tl-day").on("click", function(e){
				$("#sensor-t1-chart .x-tgl").removeClass("active");
				$(this).addClass("active");
				getStatsDataT1("day", function(d){
					if(d["incr"] == 2){
						doPlot("#sensor-t1-chart","right", "day", d["incr"], d["temp"], d["hum"]);
						$.jGrowl("Got sensor data", { 
							life: 1000, 
							closerTemplate: "<div>[ close all ]</div>",
							closeTemplate: "×" 
						});
					}
				});				
			});
			$("#sensor-t1-chart .x-tl-hour").on("click", function(e){
				$("#sensor-t1-chart .x-tgl").removeClass("active");
				$(this).addClass("active");
				getStatsDataT1("hour", function(d){
					if(d["incr"] == 2){
						doPlot("#sensor-t1-chart","right", "hour", d["incr"], d["temp"], d["hum"]);
						$.jGrowl("Got sensor data", { 
							life: 1000, 
							closerTemplate: "<div>[ close all ]</div>",
							closeTemplate: "×" 
						});
					}
				});				
			});
			$("#sensor-t2-chart .x-tl-day").on("click", function(e){
				$("#sensor-t2-chart .x-tgl").removeClass("active");
				$(this).addClass("active");
				getStatsDataT2("day", function(d){
					if(d["incr"] == 2){
						doPlot("#sensor-t2-chart","right", "day", d["incr"], d["temp"], d["hum"]);
						$.jGrowl("Got sensor data", { 
							life: 1000, 
							closerTemplate: "<div>[ close all ]</div>",
							closeTemplate: "×" 
						});
					}
				});				
			});
			$("#sensor-t2-chart .x-tl-hour").on("click", function(e){
				$("#sensor-t2-chart .x-tgl").removeClass("active");
				$(this).addClass("active");
				getStatsDataT2("hour", function(d){
					if(d["incr"] == 2){
						doPlot("#sensor-t2-chart","right", "hour", d["incr"], d["temp"], d["hum"]);
						$.jGrowl("Got sensor data", { 
							life: 1000, 
							closerTemplate: "<div>[ close all ]</div>",
							closeTemplate: "×" 
						});
					}
				});				
			});
			
			$("#svc-relay-l1-form").on("submit", function(e){
				e.preventDefault();
				//var fdata = JSON.stringify(getFormJSON($(e.target)));
				
				libUX.form.ajaxFormSubmit($(this), "/svc/relay?target=l1", "POST", function(){
					$.jGrowl("Saved relay data", { 
						life: 5000, 
						closerTemplate: "<div>[ close all ]</div>",
						closeTemplate: "×" 
					});
				});
				
			});
			$("#svc-relay-l2-form").on("submit", function(e){
				e.preventDefault();
				//var fdata = JSON.stringify(getFormJSON($(e.target)));
				
				libUX.form.ajaxFormSubmit($(this), "/svc/relay?target=l2", "POST", function(){
					$.jGrowl("Saved relay data", { 
						life: 5000, 
						closerTemplate: "<div>[ close all ]</div>",
						closeTemplate: "×" 
					});
				});
				
			});
			$("#svc-user-form").on("submit", function(e){
				e.preventDefault();
				//var fdata = JSON.stringify(getFormJSON($(e.target)));
				
				libUX.form.ajaxFormSubmit($(this), "/svc/user", "POST", function(){
					$.jGrowl("Saved account settings", { 
						life: 5000, 
						closerTemplate: "<div>[ close all ]</div>",
						closeTemplate: "×" 
					});
				});
				
			});
			
		});
			
		
	</script>
  </body>
</html>
 <!-- {{if .user_isadmin}}(user is admin){{end}} -->
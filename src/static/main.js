
var serverTimestamp = 0;
var serverOpen = false;
var myIP = "";
window.onload = function(){ 
	GetMyIp();
	LoadLogRadio();
}
var perMin = setInterval(function(){ GetTime(); }, 60000);
var perSec = setInterval(function(){ RefreshTime(); }, 1000);
var initTime = setTimeout("GetTime()", 1000);
var initDate = "2020-11-18";

function OnQuery() {
    var formData = $("#queryform").serialize();
    $.ajax(
	{   url:"/query",
	    type: "post",
	    async: true,
	    data: formData,
	    success:function(result){
		console.log(result)
		$("#result").html(result);
	    },
	    error:function(){
		$("#result").html("日志查询服务未开启");
	    }
	});
    console.log(formData)
};

function GetMyIp(){
    $.ajax(
	{   url:"/getip",
	    type: "post",
	    async: true,
	    success:function(result){
		console.log(result)
		myIP = result
		LoadSelect(); 
	    }
	});
}

function LoadSelect(){
    $.getJSON("static/server.json", function(data){
	console.log(data); 
	var optionStr = "";
	var i = 0;
	$.each(data, function (serverName){
	    console.log(serverName);
	if(serverName.indexOf("qa") != -1 || serverName.indexOf(myIP) != -1 || serverName.indexOf("全局") != -1){
	    optionStr += "<option value='" + serverName + "'";
	    if(i == 0)
		optionStr += " selected='selected'";
	    optionStr += ">";
	    optionStr += serverName;
	    optionStr += "</option>";
	    ++i;
	}
	});
	 $("#serverselect").html(optionStr);                 
    })
    console.log("select")
}

function LoadLogRadio(){
    $.getJSON("static/logpath.json", function(data){
	var radioStr = "";
	var i = 0;
	$.each(data, function (pathName){
	    console.log(pathName);
	    radioStr += "<input type='radio' name='log' value='" + data[pathName] + "'";
	    if(i == 0)
		radioStr += " checked='true'";
	    radioStr += ">";
	    radioStr += pathName + " ";
	    ++i;
	});
	$("#radios").append(radioStr);                 
    })
}

function OnAddServer(){
	$("#logsrvname").val("");
	$("#logsrvpath").val("")
	$("#model").css("display", "block");
}

function OnEnsureBtn(){
	if($("#logsrvname").val() == ""){
		alert("标识不能为空!")
		return;
	}
	if(!CheckIP($("#logsrvip").val()))
	{
		alert("IP地址不合法!");
		return;
	}
	if($("#logsrvpath").val() == ""){
		alert("日志路径不能为空!")
		return;
	}
	var formData = $("#addform").serialize();
	$.ajax(
		{   url:"/addsrv",
			type: "post",
			async: true,
			data: formData,
			success:function(result){
				console.log(result);	
				$("#serverselect").append(result);
			}
		});
	$("#model").css("display", "none");
	alert("添加成功！请查看服务器选择列表")
}

function CheckIP(value){
	var exp=/^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/;
	var reg = value.match(exp);
	if(reg == null)
		return false;
	return true;
}

function OnCancelBtn(){
	$("#model").css("display", "none");
}

function RefreshTime(){
    if(serverOpen){
	serverTimestamp++
	$("#servertime").html("服务器时间: " + TimestampToDate(serverTimestamp));
    }
}

function GetTime(){
    // console.log($("#serverselect").val())
    let formData = $('form').serialize()
    $.ajax(
	{
	    url:"/gettime",
	    type:"post",
	    async: true,
	    data: formData,
	    success:function(result){
		serverTimestamp = result;
		$("#servertime").html("服务器时间: " + TimestampToDate(serverTimestamp));
		$("#logdate").val(initDate)
		serverOpen = true
	    },
	    error:function(){
		serverOpen = false
		$("#servertime").html("远程服务未开启");
	    }
	}
    );
}

function TimestampToDate(timestamp){
    var date = new Date(timestamp * 1000);
    Y = date.getFullYear() + '/';
    initDate = date.getFullYear() + '-';
    M = (date.getMonth()+1 < 10 ? '0' + (date.getMonth()+1) : date.getMonth()+1) + '/';
    initDate += date.getMonth()+1+'-';
    D = (date.getDate() < 10 ? '0' + (date.getDate()) : date.getDate()) + ' ';
    initDate += date.getDate();
    h = (date.getHours() < 10 ? '0' + date.getHours() : date.getHours()) + ':';
    m = (date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()) + ':';
    s = date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds();
    return Y+M+D+h+m+s;
}

$('.keyclass').on('input propertychange', keyChange);

function keyChange(){
    //console.log("keyclass")
    var id = $(this).parent().attr("id");
    var type = id.indexOf("exkey") >= 0 ? "exkey" : "key"; 
    var tip = type == "exkey" ? "排除关键字" : "关键字";
    var arr = id.split(type); 
    if(arr.length != 2){
	console.log("split error");
	return;
    }
    var nextIdNum = parseInt(arr[1]) + 1;
    var nextIdStr = type + nextIdNum;
    if($(this).val() != "" && $("#"+nextIdStr).length == 0){
	var newKey = ` <div id=${nextIdStr} class="keydiv">${tip}${nextIdNum+1}: <input type="text" class="keyclass" name=${nextIdStr} >&nbsp</div>`;
	$("#"+type+"s").append(newKey);
	$(".keyclass").on('input propertychange', keyChange);
    }
    else if($(this).val() == "" && $("#"+nextIdStr).length != 0){
	var delIdNum = nextIdNum;
	while($("#"+type+delIdNum).length != 0)
	{
	    $("#"+type+delIdNum).remove();
	    ++delIdNum;
	}
    }
}


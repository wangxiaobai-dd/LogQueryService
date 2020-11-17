
var serverTimestamp = 0;
var serverOpen = false;
window.onload = function(){ LoadSelect(); LoadLogRadio(); }
var perMin = setInterval(function(){ GetTime(); }, 60000);
var perSec = setInterval(function(){ RefreshTime(); }, 1000);
var initTime = setTimeout("GetTime()", 1000);
var initDate = "2020-11-18";

function OnQuery() {
    var formData = $("form").serialize();
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

function LoadSelect(){
    $.getJSON("static/server.json", function(data){
	console.log(data); 
	var optionStr = "";
	var i = 0;
	$.each(data, function (serverName){
	    console.log(serverName);
	    optionStr += "<option value='" + serverName + "'";
	    if(i == 0)
		optionStr += " selected='selected'";
	    optionStr += ">";
	    optionStr += serverName;
	    optionStr += "</option>";
	    ++i;
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


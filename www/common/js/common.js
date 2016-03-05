/**
 * Created by Lin on 14-6-8.
 */
//子系统列表
var UrlList = [
    {url: "http://localhost:8082", name: "智翔电子实训商城",orderList:"/buyerCenter.html#/goodsMgr"},
    {url: "http://zhixiang.com", name: "电子商务商城",orderList:"/buyerCenter.html#/goodsMgr"}
];

var LinHttp = function (arg_http, arg_method, arg_url, arg_data, arg_callBack) {
    arg_http({
        method: arg_method,
        url: arg_url,
        params: arg_data
    }).success(function (response) {
        arg_callBack(response)
    }).error(function (err) {
        arg_callBack(err)
    });
};
//获得日期字符串：yyyy-MM-dd
function getDateStr() {
    var date = new Date();
    var yyyy = date.getFullYear();
    var MM = date.getMonth() + 1;
    var dd = date.getDate();
    if (MM < 10) MM = '0' + MM;
    if (dd < 10) dd = '0' + dd;
    return yyyy + "-" + MM + "-" + dd;
}
//回退
function back() {
    window.history.back();
}
//获取url中"?"符后的字串
function getUrlRequest() {
    var url = location.search; //获取url中"?"符后的字串
    //console.log(url);
    var theRequest = new Object();
    if (url.indexOf("?") != -1) {
        var str = url.substr(1);
        if (str.indexOf("&") != -1) {
            var strs = str.split("&");
            for (var i = 0; i < strs.length; i++) {
                theRequest[strs[i].split("=")[0]] = unescape(decodeURIComponent(strs[i].split("=")[1]));
            }
        } else {
            var key = str.substring(0, str.indexOf("="));
            var value = str.substr(str.indexOf("=") + 1);
            theRequest[key] = decodeURI(value);
        }
    }
    return theRequest;
}

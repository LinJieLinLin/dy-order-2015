var indexApp = "";
indexApp = angular.module("linIndex", ["ngRoute","ngResource"]);

function indexCtrl($scope,$http){
    location.href = "cart.html";
    return;
    $scope.aa = "aaa";
    var LinHttp = function (arg_method, arg_url, arg_data, arg_callBack) {
        $http({
            method: arg_method,
            url:  arg_url,
            params: arg_data
        }).success(function (response) {
            arg_callBack(response)
        }).error(function (err) {
            arg_callBack(err)
        });


//        var sendHttp = {
//            Start: function () {
//                $.ajax({
//                    type: arg_method,
//                    url: arg_url,
////                timeout:arg_timeOut,
//                    data: arg_data,
//                    dataType: "json",
//                    success: sendHttp.CallBackData,
//                    error: sendHttp.Error
//                });
//            },
//            CallBackData: function (arg_data) {
//                arg_callBack(arg_data);
//            },
//            Error: function (arg_errMsg) {
//                arg_callBack({code: 1});
//            }
//        };
//        sendHttp.Start();
    }
    var d = [{
        USER_ID :36,
        GOODS_ID:1,
        GOODS_PRICE:1,
        GOODS_NUM:1,
        GOODS_IMG :"www.baidu.com",
        GOODS_NAME:"a",
        GOODS_URL:"www.baidu.com",
        SHOP_ID :1,
        SHOP_NAME:"b",
        SHOP_URL:"www.baidu.com",
        SELLER_ID:1,
        SELLER_NAME:"c",
        FROM_URL :"www.baidu.com",
    },{
        USER_ID :2,
        GOODS_ID:2,
        GOODS_PRICE:2,
        GOODS_NUM:2,
        GOODS_IMG :"www.baidu.com",
        GOODS_NAME:"a",
        GOODS_URL:"www.baidu.com",
        SHOP_ID :1,
        SHOP_NAME:"b",
        SHOP_URL:"www.baidu.com",
        SELLER_ID:1,
        SELLER_NAME:"c",
        FROM_URL :"www.baidu.com",
    }];
    //var sb = ["1","2","3","3","2","1"];
    var sb = [{CartId:[1,2]},{CartId:[3,2]},{CartId:[4,3]}];
    var dd = angular.toJson(sb);
    var data = {data:dd,addType: "lin",age:1};
    LinHttp("POST","/l/submitOrder",data,function(arg_data){
        alert(angular.toJson(arg_data))
    })
}
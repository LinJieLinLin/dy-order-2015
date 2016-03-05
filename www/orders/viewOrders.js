/**
 * Created by LinLin on 2015/7/29.
 */
var ordersHost = "http://orders.tmp.jxzy.com/";
var ordersHost = "http://localhost:8001/";
var viewOrdersUrl = "";
function getCookie(name) {
    var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");

    if (arr = document.cookie.match(reg))

        return unescape(arr[2]);
    else
        return null;
}
function ordersListCtrl($scope, $http) {
    $scope.cb = {};
    $scope.ordersList = [];
    $scope.page = {pageNum: 1, pageSize: 3, pageSum: 1};

    $scope.cb.getOrdersList_ = function (arg_data) {
        if (!angular.isObject(arg_data) || angular.isUndefined(arg_data.code)) {
            if (arg_data) {
                console.log(arg_data);
                alert("您还未登录");
            } else {
                alert("没有数据！")
            }
            return;
        }
        if (arg_data.code == 0) {
            $scope.ordersList = [];
            $scope.page.pageSum = 1;
            if (angular.isDefined(arg_data.data.page) && angular.isDefined(arg_data.data.ordersList) && angular.isArray(arg_data.data.ordersList)) {
                $scope.page.pageSum = arg_data.data.page;
                $scope.ordersList = arg_data.data.ordersList;
                console.log($scope.ordersList);
                var temL = $scope.ordersList.length;
                if (UrlList && UrlList.length > 0 && angular.isArray($scope.ordersList) || temL > 0) {
                    for (var i = 0; i < temL; i++) {
                        var temUrlListL = UrlList.length;
                        for (var k = 0; k < temUrlListL; k++) {
                            if ($scope.ordersList[i].fromUrl.indexOf(UrlList[k].url) > -1) {
                                $scope.ordersList[i].fromUrlName = UrlList[k].name;
                                $scope.ordersList[i].fromUrl = UrlList[k].url;
                                break;
                            }
                        }
                    }
                }
            } else {
                console.log("暂无数据");
            }
        } else {
            if (angular.isDefined(arg_data.msg)) {
                alert(arg_data.msg);
            }
        }
    };
    $scope.getOrdersList = function (arg_pNum, arg_pSize, arg_status) {
        if (angular.isUndefined(arg_pNum) || angular.isUndefined(arg_pSize) || angular.isUndefined(arg_status)) {
            return;
        }

        var data = {
            PNum: arg_pNum,
            PSize: arg_pSize,
            Status: arg_status
            //-1为所有状态
        };
        var token = getCookie("token");
        if (angular.isUndefined(token) || token == "") {
            alert("您还未登录");
        }
        viewOrdersUrl = "l/getOrdersData";
        viewOrdersUrl += "?token=" + token;
        data = angular.toJson(data);
        LinHttp($http, "GET", ordersHost + viewOrdersUrl, {data: data}, $scope.cb.getOrdersList_)
    };
    $scope.changePage = function (arg_type) {
        if (arg_type != -1 && arg_type != 1) {
            return
        }
        $scope.page.pageNum = $scope.page.pageNum + arg_type;
        if ($scope.page.pageNum < 1) {
            $scope.page.pageNum = $scope.page.pageNum - arg_type;
            return;
        } else if ($scope.page.pageNum > $scope.page.pageSum) {
            $scope.page.pageNum = $scope.page.pageNum - arg_type;
            return;
        }
        $scope.getOrdersList($scope.page.pageNum, $scope.page.pageSize, -1);
    };
    $scope.getOrdersList($scope.page.pageNum, $scope.page.pageSize, -1);
}
var ordersApp = angular.module("ordersApp", ["ngResource"]);

function orderCtrl($scope, $http) {
    $scope.cartData = angular.fromJson(localStorage.CARTDATA);
    $scope.selectData = angular.fromJson(localStorage.PRICE);
    $scope.cb = {};
    $scope.cb.getCart_ = function (arg_data) {
        if (!angular.isObject(arg_data) || angular.isUndefined(arg_data.code)) {
            if (arg_data) {
                console.log(arg_data);
                alert("您还未登录")
            } else {
                alert("没有数据！")
            }
            return;
        }
        if (arg_data.code == 0) {
            $scope.cartData = [];
            if (!angular.isArray(arg_data.data)) {
                console.log("购物车为空");
                return;
            }
            //var temL = arg_data.data.length;
            //for (var i = 0; i < temL; i++) {
            //    $scope.cartData = $scope.cartData.concat(arg_data.data[i].sellers);
            //}
            $scope.cartData = arg_data.data;
            var temL = $scope.cartData.length;
            $scope.selectData = {count: 0, price: 0};
            for (var i = 0; i < temL; i++) {
                $scope.cartData[i].haveSelect = true;
                var temGoodsL = $scope.cartData[i].goods.length;
                var temCount = 0;
                var temPrice = 0;
                var temUrlListL = UrlList.length;
                for (var k = 0; k < temUrlListL; k++) {
                    if ($scope.cartData[i].fromUrl.indexOf(UrlList[k].url) > -1) {
                        $scope.cartData[i].fromUrlName = UrlList[k].name;
                        $scope.cartData[i].fromUrl = UrlList[k].url;
                        break;
                    }
                }
                for (var j = 0; j < temGoodsL; j++) {
                    $scope.cartData[i].goods[j].isSelected = true;
                    $scope.selectData.count++;
                    temCount += $scope.cartData[i].goods[j].goodsNum;
                    temPrice += $scope.cartData[i].goods[j].goodsPrice * $scope.cartData[i].goods[j].goodsNum;
                }
                $scope.cartData[i].count = temCount;
                $scope.cartData[i].price = temPrice;
                $scope.selectData.price += temPrice;
            }
            localStorage.CARTDATA = angular.toJson($scope.cartData);
            localStorage.PRICE = angular.toJson($scope.selectData);
            //$scope.$apply(function(){
            //    $scope.cartData = $scope.cartData;
            //});
        } else {
            console.log(arg_data.msg);
            alert("您还未登录")
        }
    };

    var urlData = getUrlRequest();
    if (angular.isDefined(urlData.cartId)) {
        console.log(urlData.cartId);
        var cartId = angular.fromJson(urlData.cartId);
        if (!angular.isArray(cartId)) {
            alert("商品数据有误");
        }
        $scope.cartData = [];
        $scope.selectData = {};
        $scope.getCart = function () {
            var data = {
                CartId: cartId
            };
            data = angular.toJson(data);
            LinHttp($http, "GET", "/l/getCart", {data: data}, $scope.cb.getCart_)
        };
        $scope.getCart();
    } else {
        if (angular.isUndefined($scope.cartData) || !angular.isArray($scope.cartData)) {
            console.log("没有商品数据！");
            return;
        }
        if (angular.isUndefined($scope.selectData) || !angular.isObject($scope.selectData)) {
            console.log("没有数据！");
            return;
        }
        temL = $scope.cartData.length;
        for (var i = 0; i < temL; i++) {
            //$scope.cartData[i].haveSelect = true;
            var temGoodsL = $scope.cartData[i].goods.length;
            var temCount = 0;
            var temPrice = 0;
            var temUrlListL = UrlList.length;
            for (var k = 0; k < temUrlListL; k++) {
                if ($scope.cartData[i].fromUrl.indexOf(UrlList[k].url) > -1) {
                    $scope.cartData[i].fromUrlName = UrlList[k].name;
                    $scope.cartData[i].fromUrl = UrlList[k].url;
                    break;
                }
            }
            for (var j = 0; j < temGoodsL; j++) {
                //$scope.cartData[i].goods[j].isSelected = true;
                temCount += $scope.cartData[i].goods[j].goodsNum;
                temPrice += $scope.cartData[i].goods[j].goodsPrice * $scope.cartData[i].goods[j].goodsNum;
            }
            $scope.cartData[i].count = temCount;
            $scope.cartData[i].price = temPrice;
        }
    }
    $scope.submitOrder = function () {
        var ids = [];
        var temL = $scope.cartData.length;
        for (var i = 0; i < temL; i++) {
            var temGoodsL = $scope.cartData[i].goods.length;
            var CartIds = [];
            var isPush = false;
            for (var j = 0; j < temGoodsL; j++) {
                if ($scope.cartData[i].goods[j].isSelected) {
                    isPush = true;
                    CartIds.push($scope.cartData[i].goods[j].id);
                }
            }
            if (isPush) {
                ids.push({CartId: CartIds})
            }
        }
        if (ids.length < 1) {
            alert("商品数据有误！");
            return
        }

        //console.log(ids);
        LinHttp($http, "POST", "/l/submitOrder", {data: angular.toJson(ids)}, $scope.cb.submitOrder_)
    };
    console.log($scope.cartData);
    $scope.cb.submitOrder_ = function (arg_data) {
        if (!angular.isObject(arg_data) || angular.isUndefined(arg_data.code)) {
            if (arg_data) {
                alert(arg_data)
            } else {
                alert("没有数据！")
            }
            return;
        }
        if (arg_data.code == 0 && angular.isDefined(arg_data.data)) {
            if (arg_data.data.orderNo=="free"){
                location.href = "/paymentStatus.html?type=success&orderNos=free";
                return
            }
            localStorage.OrdersNo = angular.toJson(arg_data.data);
            location.href = "/pay.html";
            return
        } else {
            alert(arg_data.msg);
        }
    };
}
function payCtrl($scope, $http) {
    $scope.isSuccess = false;
    $scope.msg = "";
    var urlData = getUrlRequest();
    if (angular.isObject(urlData) && angular.isDefined(urlData.type)) {
        if (urlData.type == "fail") {
            $scope.isSuccess = false;
            $scope.msg = urlData.msg
        } else if (urlData.type == "success") {
            $scope.isSuccess = true;
        }
    }
    $scope.cartData = angular.fromJson(localStorage.CARTDATA);
    $scope.selectData = angular.fromJson(localStorage.PRICE);
    $scope.ordersNo = angular.fromJson(localStorage.OrdersNo);
    $scope.cb = {};
    $scope.goodsName = [];
    $scope.time = getDateStr();

    if (angular.isUndefined($scope.cartData) || !angular.isArray($scope.cartData)) {
        console.log("没有商品数据！");
        return;
    }
    if (angular.isUndefined($scope.selectData) || !angular.isObject($scope.selectData)) {
        console.log("没有数据！");
        return;
    }
    if (angular.isUndefined($scope.ordersNo) || !angular.isObject($scope.ordersNo)) {
        console.log("没有数据！");
        return;
    }

    $scope.sonOrdersList = [];
    //if($scope.cartData.length){
    //    $scope.sonOrdersList = $scope.cartData[0].fromUrl;
    //}
    $scope.getGoodsName = function () {
        $scope.goodsName = [];
        var temL = $scope.cartData.length;
        for (var i = 0; i < temL; i++) {
            var temUrlListL = UrlList.length;
            for (var k = 0; k < temUrlListL; k++) {
                if ($scope.cartData[i].fromUrl.indexOf(UrlList[k].url) > -1) {
                    var temAdd = true;
                    $scope.sonOrdersList.forEach(function (e) {
                        if (e.url == UrlList[k].url) {
                            temAdd = false;
                        }
                    });
                    if (temAdd) {
                        $scope.sonOrdersList.push(UrlList[k]);
                    }
                    break;
                }
            }
            var temGoodsL = $scope.cartData[i].goods.length;
            for (var j = 0; j < temGoodsL; j++) {
                if ($scope.cartData[i].goods[j].isSelected) {
                    $scope.goodsName.push($scope.cartData[i].goods[j].goodsName);
                }
            }
        }
        console.log($scope.sonOrdersList)
    };
    $scope.getGoodsName();
    $scope.pay = function () {
        console.log($scope.ordersNo);
        //LinHttp($http, "POST", "/l/alipay", {publicOrderNo: $scope.ordersNo.publicOrderNo}, function () {});
        window.open("/l/alipay?orderNo=" + $scope.ordersNo.orderNo);
    };
}
function viewOrdersCtrl($scope, $http) {
    $scope.viewOrders = "orders/ordersList.html";
}
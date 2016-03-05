var cartApp = angular.module("cartApp", ["ngResource"]);

function cartCtrl($scope, $http) {
    $scope.cartData = [];
    $scope.cb = {};
    $scope.isSelectAll = false;
    $scope.selectData = {count: 0, price: 0};
    $scope.dataCount = 0;
    
    $scope.checkSelectAll = function () {
        var temL = $scope.cartData.length;
        for (var i = 0; i < temL; i++) {
            var temGoodsL = $scope.cartData[i].goods.length;
            $scope.cartData[i].isSelected = $scope.isSelectAll;
            for (var j = 0; j < temGoodsL; j++) {
                if (!$scope.cartData[i].haveSelect) {
                    $scope.cartData[i].haveSelect = true;
                }
                $scope.cartData[i].goods[j].isSelected = $scope.isSelectAll;
            }
        }
        $scope.countPrice()
    };
    $scope.checkSelectStoreAll = function (arg_data) {
        console.log(arg_data);
        if (!arg_data.haveSelect) {
            arg_data.haveSelect = true;
        }
        var temL = arg_data.goods.length;
        for (var i = 0; i < temL; i++) {
            arg_data.goods[i].isSelected = arg_data.isSelected;
            console.log(arg_data.goods[i]);
        }
        $scope.countPrice()
    };

    $scope.countPrice = function () {
        $scope.selectData = {count: 0, price: 0};
        var temL = $scope.cartData.length;
        for (var i = 0; i < temL; i++) {
            var temGoodsL = $scope.cartData[i].goods.length;
            for (var j = 0; j < temGoodsL; j++) {
                if ($scope.cartData[i].goods[j].isSelected) {
                    if (!$scope.cartData[i].haveSelect) {
                        $scope.cartData[i].haveSelect = true;
                    }
                    $scope.selectData.count=$scope.selectData.count+$scope.cartData[i].goods[j].goodsNum;
                    $scope.selectData.price = $scope.selectData.price + $scope.cartData[i].goods[j].goodsPrice * $scope.cartData[i].goods[j].goodsNum;
                }
            }
        }
    };


    $scope.getCart = function () {
        LinHttp($http, "GET", "/l/getCart", {data: {}}, $scope.cb.getCart_)
    };
    $scope.rmGoods = function (arg_data) {
        if (angular.isUndefined(arg_data.id)) {
            return;
        }
        var editCartData = [{CartId: arg_data.id, Num: 1, Status: -1}];
        if ($scope.dataCount > 0) {
            $scope.dataCount = $scope.dataCount - 1;
        }
        $scope.editCart(editCartData);
    };
    $scope.rmAllGoods = function () {
        //暂时不做
    };
    $scope.editCart = function (arg_data) {
        var data = {data: angular.toJson(arg_data)};
        LinHttp($http, "POST", "/l/editCart", data, $scope.cb.editCart_)
    };

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
            if (angular.isArray(arg_data.data)) {
                var temL = arg_data.data.length;
                $scope.dataCount = 0;
                //for (var i = 0; i < temL; i++) {
                //    $scope.cartData = $scope.cartData.concat(arg_data.data[i].sellers);
                //}
                $scope.cartData = arg_data.data;
                temL = $scope.cartData.length;
                for (var i = 0; i < temL; i++) {
                    var temUrlListL = UrlList.length;
                    for(var k=0;k<temUrlListL;k++){
                        if($scope.cartData[i].fromUrl.indexOf(UrlList[k].url)>-1){
                            $scope.cartData[i].fromUrlName = UrlList[k].name;
                            $scope.cartData[i].fromUrl = UrlList[k].url;
                            break;
                        }
                    }
                    var temGoodsL = $scope.cartData[i].goods.length;
                    for (var j = 0; j < temGoodsL; j++) {
                        $scope.dataCount++;
                    }
                }
            }
        } else {
            console.log(arg_data.msg);
            alert("您还未登录")
        }
    };
    $scope.cb.editCart_ = function (arg_data) {
        if (!angular.isObject(arg_data) || angular.isUndefined(arg_data.code)) {
            if (arg_data) {
                console.log(arg_data);
                alert("您还未登录")
            } else {
                alert("没有数据！")
            }
            return
        }
        if (arg_data.code == 0) {
            $scope.getCart();
        } else {
            console.log(arg_data.msg);
            alert("您还未登录");
        }
    };

    $scope.goToOrders = function () {
        if ($scope.selectData.count < 1) {
            return;
        }
        localStorage.CARTDATA = angular.toJson($scope.cartData);
        localStorage.PRICE = angular.toJson($scope.selectData);
        location.href = "orders.html";
    };
    $scope.getCart();
}
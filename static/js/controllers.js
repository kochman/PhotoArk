var app = angular.module('app', ['angularLazyImg', 'ngAnimate']);

app.controller('HomeCtrl', function($scope, $http, $sce) {
	$http.get('api/filters').success(function(data) {
		$scope.photos = [];
		$scope.filters = data;

		for (var key in $scope.filters) {
			if (!$scope.filters.hasOwnProperty(key)) {
				continue;
			}
			var filters = $scope.filters[key];
			for (var i = 0; i < filters.length; i++) {
				filters[i] = {name: filters[i], enabled: false};
			}
		}
	});

	$scope.filterChanged = function() {
		var filterParams = {};
		for (var key in $scope.filters) {
			if (!$scope.filters.hasOwnProperty(key)) {
				continue;
			}
			var filters = $scope.filters[key];
			filterParams[key] = [];
			for (var i = 0; i < filters.length; i++) {
				if (filters[i].enabled) {
					filterParams[key].push(filters[i].name);
				}
			}
		}

		$http.get('api/filter', {params: filterParams}).success(function(data) {
			$scope.photos = data;
		});
	}

	$scope.photoDetail = function(photo) {
		$scope.photoDetailModal = true;
		$scope.photoDetailModalPhoto = photo;
		$http.get('api/metadata', {params: {photo: photo}}).success(function(data) {
			$scope.photoDetailModalMetadata = data;
		});
	}

	$scope.closePhotoDetail = function() {
		$scope.photoDetailModal = false;
		delete $scope.photoDetailModalPhoto;
		delete $scope.photoDetailModalMetadata;
	}
});

var app = angular.module('app', ['angularLazyImg', 'ngAnimate']);

app.config(function($locationProvider) {
	$locationProvider.html5Mode(true);
});

app.controller('HomeCtrl', function($scope, $http, $location) {
	$http.get('api/filters').success(function(data) {
		$scope.photos = [];
		$scope.filters = data;

		var splitPath = $location.path().split('/');
		splitPath.shift();  // ignore first because it is empty
		var activeURLFilters = {}
		for (var i = 0; i < splitPath.length; ++i) {
			var split = splitPath[i].split(':');
			var key = split[0];
			var val = split[1];
			if (!activeURLFilters[key]) {
				activeURLFilters[key] = [];
			}
			activeURLFilters[key].push(val);
		}

		for (var key in $scope.filters) {
			if (!$scope.filters.hasOwnProperty(key)) {
				continue;
			}
			var filters = $scope.filters[key];
			for (var i = 0; i < filters.length; i++) {
				filters[i] = {name: filters[i], enabled: false};
				if (activeURLFilters[key]) {
					var filter = filters[i];
					if (activeURLFilters[key].indexOf(filter.name) != -1) {
						filter.enabled = true;
					}
				}
			}
		}
		$scope.filterChanged();
		$scope.initialized = true;
	});

	$scope.filterChanged = function() {
		var filterParams = {};
		var newPath = '';
		for (var key in $scope.filters) {
			if (!$scope.filters.hasOwnProperty(key)) {
				continue;
			}
			var filters = $scope.filters[key];
			filterParams[key] = [];
			for (var i = 0; i < filters.length; i++) {
				if (filters[i].enabled) {
					filterParams[key].push(filters[i].name);
					newPath += '/' + key + ':' + filters[i].name
				}
			}
		}

		$location.path(newPath);

		// check if filters are not empty
		if (newPath != '') {
			$http.get('api/filter', {params: filterParams}).success(function(data) {
				$scope.photos = data;
			});
		} else {
			$scope.photos = [];
		}
	}

	$scope.locationChanged = function() {
		var splitPath = $location.path().split('/');
		splitPath.shift();  // ignore first because it is empty
		var activeURLFilters = {}
		for (var i = 0; i < splitPath.length; ++i) {
			var split = splitPath[i].split(':');
			var key = split[0];
			var val = split[1];
			if (!activeURLFilters[key]) {
				activeURLFilters[key] = [];
			}
			activeURLFilters[key].push(val);
		}

		for (var key in $scope.filters) {
			if (!$scope.filters.hasOwnProperty(key)) {
				continue;
			}
			var filters = $scope.filters[key];
			for (var i = 0; i < filters.length; i++) {
				var filter = filters[i];
				if (activeURLFilters[key]) {
					if (activeURLFilters[key].indexOf(filter.name) != -1) {
						filter.enabled = true;
					} else {
						filter.enabled = false;
					}
				} else {
					filter.enabled = false;
				}
			}
		}

		$scope.filterChanged();
	}

	$scope.$on('$locationChangeSuccess', function() {
		if ($scope.initialized) {
			$scope.locationChanged();
		}
	});

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

(function () {

	/**
	 * The PackageService provides access to Dart package
	 * data using REST API calls.
	 */
	function PackageService($http) {
		this.getRecent = function() {
			return $http.get('/api/packages');
		};
		this.get = function(name) {
			return $http.get('/api/packages/'+name);
		};
	}

	angular
		.module('app')
		.service('packages', PackageService);
})();
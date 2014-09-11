(function () {

	/**
	 * The PackageService provides access to Dart package
	 * data using REST API calls.
	 */
	function PackageService($http) {
		this.getRecent = function() {
			return $http.get('https://storage.googleapis.com/brad_dart_test/_recent.json');
		};
		this.get = function(name) {
			return $http.get('https://storage.googleapis.com/brad_dart_test/'+name+'/package.json');
		};
	}

	angular
		.module('app')
		.service('packages', PackageService);
})();
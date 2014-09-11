(function () {

	/**
	 * The BuildService provides access to package build
	 * data using REST API calls.
	 */
	function BuildService($http) {
		this.get = function(name, version) {
			return $http.get('https://storage.googleapis.com/brad_dart_test/'+name+'/'+version+'/results.json');
		};
		this.getOutput = function(name, version) {
			return $http.get('https://storage.googleapis.com/brad_dart_test/'+name+'/'+version+'/output.txt');
		};
	}

	angular
		.module('app')
		.service('builds', BuildService);
})();
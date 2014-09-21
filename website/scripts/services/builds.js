(function () {

	/**
	 * The BuildService provides access to package build
	 * data using REST API calls.
	 */
	function BuildService($http) {
		this.get = function(name, version, channel, sdk) {
			return $http.get('/api/packages/'+name+'/'+version+'/channel/'+(channel || 'stable')+'/sdk/'+(sdk || 'latest'));
		};
		this.getOutput = function(name, version, channel, sdk) {
			return $http.get('/api/packages/'+name+'/'+version+'/channel/'+(channel || 'stable')+'/sdk/'+sdk+'/stdout.txt');
		};
	}

	angular
		.module('app')
		.service('builds', BuildService);
})();
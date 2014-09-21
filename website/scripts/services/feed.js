(function () {

	/**
	 * The FeedService provides access to the build
	 * feed using REST API calls.
	 */
	function FeedService($http) {
		this.get = function() {
			return $http.get('/api/feed');
		};
	}

	angular
		.module('app')
		.service('feed', FeedService);
})();
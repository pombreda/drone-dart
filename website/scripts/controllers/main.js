(function () {

	/**
	 * Main controller responsible for displaying
	 * the main homepage and rendering recent package
	 * uploads and builds.
	 */	
	function MainCtrl($scope, $interval, feed) {
		feed.get().then(function(recent){
			$scope.recent = recent.data;
		}).catch(function(err){
			$scope.error = err;
		});

		// check for new data every 5 minutes
		//$interval(function() {
		//	feed.get().then(function(recent){
		//		$scope.recent = recent.data;
		//	});
		//}, 60000);
	}

	angular
		.module('app')
		.controller('MainCtrl', MainCtrl);
})();
(function () {

	/**
	 * Build controller responsible for displaying
	 * the build page and rendering the build results.
	 */	
	function BuildCtrl($scope, $routeParams, builds) {
		var package = $routeParams.package;
		var version = $routeParams.version;
		var channel = $routeParams.channel;
		var sdk     = $routeParams.sdk;

		// gets the build data from the server
		builds.get(package, version, channel, sdk).then(function(build){
			$scope.build = build.data;
			channel = build.data.channel;
			sdk     = build.data.sdk;

			// gets the build output from the server
			builds.getOutput(package, version, channel, sdk).then(function(output){
				$scope.output = output.data;
			}).catch(function(error){
				$scope.error = error;
			});

		}).catch(function(error){
			$scope.error = error;
		});
	}

	angular
		.module('app')
		.controller('BuildCtrl', BuildCtrl);
})();
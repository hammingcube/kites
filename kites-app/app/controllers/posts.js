import Ember from 'ember';

export default Ember.ArrayController.extend({
	sortProperties: ['timestamp'],
	sortAscending: false,
	actions: {
		publishPost: function() {
			var newPost = this.store.createRecord('post', {
				title: this.get('title'),
				body: this.get('body'),
				timestamp: new Date().getTime()
			});
			newPost.save();
		}
	}
});

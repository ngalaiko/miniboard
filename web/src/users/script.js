import './components/UserFeedController/UserFeedController.js'

document.querySelector('#user-feed-controller').addEventListener('FeedClicked', (e) => {
    console.log(e.detail.id, "clicked")
})

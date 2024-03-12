const uploadButton = document.getElementById('upload-button');
const uploadInput = document.getElementById('upload');
const videoTitleInput = document.getElementById('video-title');
const uploadVideoButton = document.getElementById('upload-video');
const videoList = document.getElementById('video-list');

let isUploadVisible = false;

uploadButton.addEventListener('click', () => {
    if (!isUploadVisible) {
        uploadInput.style.display = 'block';
        videoTitleInput.style.display = 'block';
        uploadVideoButton.style.display = 'block';
        isUploadVisible = true;
    } else {
        uploadInput.style.display = 'none';
        videoTitleInput.style.display = 'none';
        uploadVideoButton.style.display = 'none';
        isUploadVisible = false;
    }
});

uploadVideoButton.addEventListener('click', () => {
    const file = uploadInput.files[0];
    if (file) {
        const formData = new FormData();
        formData.append('video', file);
        formData.append('title', videoTitleInput.value);

        fetch('/upload', {
            method: 'POST',
            body: formData
        })
            .then(response => response.json())
            .then(data => {
                console.log(data);
                uploadInput.style.display = 'none';
                videoTitleInput.style.display = 'none';
                uploadVideoButton.style.display = 'none';
                uploadButton.style.display = 'block'; // 显示"上传视频"按钮
                loadVideoList(); // 调用加载视频列表的函数
            });
    }
});

// 加载视频列表的函数
// 修改加载视频列表的函数
function loadVideoList() {
    fetch('/videos')
        .then(response => response.json())
        .then(data => {
            videoList.innerHTML = '';

            data.forEach(video => {
                const li = document.createElement('li');
                const videoElement = document.createElement('video');
                videoElement.src = video.url;
                videoElement.controls = true;

                const titleLink = document.createElement('a');
                titleLink.href = '/videoplay/' + video.url.split('/').pop();
                titleLink.textContent = video.title;

                // 创建一个容器元素，用于包装视频和标题
                const videoContainer = document.createElement('div');
                videoContainer.classList.add('video-container');
                videoContainer.appendChild(videoElement);
                videoContainer.appendChild(titleLink);

                li.appendChild(videoContainer);
                videoList.appendChild(li);
            });


        });
}



// 初始加载视频列表
loadVideoList();

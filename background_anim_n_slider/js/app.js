import Swiper from 'swiper/bundle';
import 'swiper/css/bundle';


const video = document.querySelector('.background-video')

const swiperText = new Swiper('.swiper', {
    speed: 800,
    loop: true,
    // effect: 'slide',
    effect: 'fade',
    mousewheel: {},
    keyboard: {},
    pagination: {
        el: '.swiper-pagination',
        clickable: true,
    },
    navigation: {
        prevEl: '.swiper-button-prev',
        nextEl: '.swiper-button-next',
    },
})

swiperText.on('slideChange', function() {
    gsap.to(video, 1.6, {
        currentTime: (video.duration / (this.slides.length - 1)) * this.realIndex,
        ease: Power4.easeOut,
    })
})

let timeoutId;

swiperText.on('slideChangeTransitionStart', function() {
    if (timeoutId) {
        clearTimeout(timeoutId);
    }
    video.classList.add('change')
})

swiperText.on('slideChangeTransitionEnd', function() {
    video.classList.remove('change')
})

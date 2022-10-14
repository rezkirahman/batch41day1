let dataBlog = []

function addBlog(event) {
    event.preventDefault()

    let title = document.getElementById("input-title").value
    let startDate = document.getElementById("input-startdate").value
    let endDate = document.getElementById("input-enddate").value
    let description = document.getElementById("input-content").value
    let tech = document.getElementById("input-tech").value
    let image = document.getElementById("input-blog-image").files[0]

    // buat url gambar nantinya tampil
    image = URL.createObjectURL(image)
    console.log(image)

    let blog = {
        title,
        startDate,
        endDate,
        tech,
        description,
        image,
        postAt: new Date(),
        author: "Rezki Rahman"
    }

    dataBlog.push(blog)
    console.log(dataBlog)

    renderBlog()
}

function renderBlog() {
    document.getElementById("contents").innerHTML = ''

    for (let index = 0; index < dataBlog.length; index++) {
        console.log("test",dataBlog[index])

        document.getElementById("contents").innerHTML += `
        <div class="blog-list-item">
            <div class="blog-image">
                <img src="${dataBlog[index].image}">
            </div>
            <div class="blog-content">
                <h3>
                    <a href="blog-detail.html" target="_blank" style:"color: black;">
                        ${dataBlog[index].title}
                    </a>
                </h3>
                <div class="detail-blog-content">
                    ${getFullTime(dataBlog[index].postAt)} | ${dataBlog[index].author}                    
                </div>
                <p> 
                    ${getDistanceTime(dataBlog[index].postAt)}
                </p>
                <p>
                    ${dataBlog[index].description}
                </p>
                <p>
                    ${dataBlog[index].tech}
                </p>
            
                <div style="text-align: left; font-size: 20px; " >
                    <a href="" style="color: black;"><i class="fa-brands fa-instagram"></i></a>
                    <a href="" style="color: black;"><i class="fa-brands fa-facebook" ></i></a>
                    <a href="" style="color: black;"><i class="fa-brands fa-twitter"></i></a>
                    <a href="" style="color: black;"><i class="fa-brands fa-linkedin"></i></a>
                </div>
            </div>
            <div class="btn-group">
                    <button class="btn-detail" style="margin-right: 4px;">Edit Post</button>
                    <button class="btn-detail" style="margin-left: 4px;">Post Blog</button>
            </div>
        </div>
        `
    }
}

function getFullTime(time) {
    //time = new Date()
    //console.log(time)
    let monthName = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dev']

    let date = time.getDate()
    console.log(date)

    let monthIndex = time.getMonth()
    console.log(monthIndex)

    let year = time.getFullYear()
    console.log(year)

    let hours = time.getHours()
    let minutes = time.getMinutes()

    if (hours <=9) {
        hours = "0" + hours 
    } else if (minutes<=9) {
        minutes = "0" + minutes
    }

    return `${date} ${monthName[monthIndex]} ${year} ${hours}:${minutes} WIB` 

}

function getDistanceTime(time) {
    let timeNow = new Date()
    let timePost = time
    let distance = timeNow - timePost //milisecond
    console.log(distance)

    let milisecond = 1000
    let secondInHours = 3600
    let hoursInDay = 24

    let distanceDay = Math.floor(distance/ (milisecond * secondInHours * hoursInDay))
    let distancehours = Math.floor(distance / (milisecond * 60 * 60))
    let distanceMinutes = Math.floor(distance / (milisecond * 60))
    let distanceSecond = Math.floor(distance / milisecond)

    if (distanceDay > 0){
        return `${distanceDay} day ago`
    } else if (distancehours > 0 ){
        return `${distancehours} hour(s) ago`
    } else if (distanceMinutes > 0){
        return `${distanceMinutes} minute(s) ago`
    } else {
        return `${distanceSecond} second ago`
    }
}

setInterval(function() {
    renderBlog()
}, 2000)

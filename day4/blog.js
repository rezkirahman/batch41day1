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
                    <a href="blog-detail.html" target="_blank">
                        ${dataBlog[index].title}
                    </a>
                </h3>
                <div class="detail-blog-content">
                    12 Jul 2021 22:30 WIB | Rezki Rahman
                </div>
                <p>
                    ${dataBlog[index].content}
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
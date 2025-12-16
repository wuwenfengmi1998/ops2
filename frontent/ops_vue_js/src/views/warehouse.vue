<template>
  <div class="page-wrapper">
    <!-- Page header -->
    <div class="page-header d-print-none">
      <div class="container-xl">
        <div class="row g-2 align-items-center">
          <div class="col">
            <h2 class="page-title">Cards</h2>
          </div>
        </div>
      </div>
    </div>
    <!-- Page body -->
    <div class="page-body">
      <div class="container-xl">
        <div class="row row-cards">
          <div class="col-md-6 col-lg-3">
            <div class="card">
              <!-- Photo -->
              <div
                class="img-responsive img-responsive-21x9 card-img-top"
                style="
                  background-image: url(./static/photos/home-office-desk-with-macbook-iphone-calendar-watch-and-organizer.jpg);
                "
              ></div>
              <div class="card-body">
                <h3 class="card-title">Card with top image</h3>
                <p class="text-secondary card_text">
                  Lorem ipsum dolor sit amet, consectetur adipisicing elit.
                 
                </p>
              </div>
            </div>
          </div>
          <div class="col-md-6 col-lg-3">
            <div class="card">
              <!-- Photo -->
              <div
                class="img-responsive img-responsive-21x9 card-img-top"
                style="
                  background-image: url(./static/photos/home-office-desk-with-macbook-iphone-calendar-watch-and-organizer.jpg);
                "
              ></div>
              <div class="card-body">
                <h3 class="card-title">Card with top image</h3>
                <p class="text-secondary card_text">
                  Lorem ipsum dolor sit amet, consectetur adipisicing elit.
                  Aperiam deleniti fugit incidunt, iste, itaque minima neque
                  pariatur perferendis sed suscipit velit vitae voluptatem. Lorem, ipsum dolor sit amet consectetur adipisicing elit. Illo optio et fuga omnis ipsa, odit repellendus iste doloremque est, nam eius quisquam perspiciatis deserunt. Quasi tempore velit architecto corporis voluptatibus!
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


<style scoped>

.card_text{
   display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2; /* 限制显示4行 */
  overflow: hidden;
  line-height: 1.5;
  min-height: calc(1.5em * 2); /* 最小高度 */
}

</style>

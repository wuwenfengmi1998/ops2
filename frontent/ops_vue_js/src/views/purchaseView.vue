<template>
   
  <div class="col-12">
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">Invoices</h3>
      </div>
      <div class="card-body border-bottom py-3">
        <div class="d-flex">
          <div class="text-secondary">
            Show
            <div class="mx-2 d-inline-block">
              <input
                type="text"
                class="form-control form-control-sm"
                value="8"
                size="3"
                aria-label="Invoices count"
              />
            </div>
            entries
          </div>
          <div class="ms-auto text-secondary">
            Search:
            <div class="ms-2 d-inline-block">
              <input
                type="text"
                class="form-control form-control-sm"
                aria-label="Search invoice"
              />
            </div>
          </div>
        </div>
      </div>
      <div class="table-responsive">
        <table class="table card-table table-vcenter text-nowrap datatable">
          <thead>
            <tr>
              <th class="w-1">
                <input
                  class="form-check-input m-0 align-middle"
                  type="checkbox"
                  aria-label="Select all invoices"
                />
              </th>
              <th class="col-1">
                No.
                <!-- Download SVG icon from http://tabler-icons.io/i/chevron-up -->
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="icon icon-sm icon-thick"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  stroke-width="2"
                  stroke="currentColor"
                  fill="none"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                  <path d="M6 15l6 -6l6 6" />
                </svg>
              </th>
              <th class="col-5">物品名称</th>
              <th class="col-1">用途</th>
              <th class="w-1">单位</th>
              <th class="w-1">数量</th>
              <th class="w-1">单价</th>
              <th class="w-1">总价</th>
              <th class="w-1">创建日期</th>
              <th class="w-1">更新日期</th>
              <th class="w-1">最新状态</th>
              <th class="w-1"></th>
            </tr>
          </thead>
          <tbody>
            <tr>
                <td>
                    <input
                    class="form-check-input m-0 align-middle"
                    type="checkbox"
                    aria-label="Select invoice"
                    />
                </td>
                <td><span class="text-muted">001</span></td>
                <td>
                    办公室用纸
                </td>
                <td>办公用品</td>
                <td>包</td>
                <td>10</td>
                <td>15.00</td>
                <td>150.00</td>
                <td>2024-06-01</td>
                <td>2024-06-05</td>
                <td>
                    <span class="badge bg-success me-1"></span> 已完成
                </td>
                <td class="text-end">
                   
                </td>
            </tr>
           
          </tbody>
        </table>
      </div>
      <div class="card-footer d-flex align-items-center">
        <p class="m-0 text-secondary">
          Showing <span>1</span> to <span>8</span> of <span>16</span> entries
        </p>
        <ul class="pagination m-0 ms-auto">
          <li class="page-item disabled">
            <a class="page-link" href="#" tabindex="-1" aria-disabled="true">
              <!-- Download SVG icon from http://tabler-icons.io/i/chevron-left -->
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="icon"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                stroke-width="2"
                stroke="currentColor"
                fill="none"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M15 6l-6 6l6 6" />
              </svg>
              prev
            </a>
          </li>
          <li class="page-item"><a class="page-link" href="#">1</a></li>
          <li class="page-item active"><a class="page-link" href="#">2</a></li>
          <li class="page-item"><a class="page-link" href="#">3</a></li>
          <li class="page-item"><a class="page-link" href="#">4</a></li>
          <li class="page-item"><a class="page-link" href="#">5</a></li>
          <li class="page-item">
            <a class="page-link" href="#">
              next
              <!-- Download SVG icon from http://tabler-icons.io/i/chevron-right -->
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="icon"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                stroke-width="2"
                stroke="currentColor"
                fill="none"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M9 6l6 6l-6 6" />
              </svg>
            </a>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

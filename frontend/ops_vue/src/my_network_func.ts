import axios from 'axios'
import { myfuncs } from './myfunc'

interface re_data {
  statusCode: number
  cookie?: string
  error?: string
  data?: JSON
}

export const my_network_func = {
  post_json(path: string, json: JSON, callback: (i: re_data) => void) {
    const head_path = '/api/v1'
    const url = head_path + path
    const cookie = myfuncs.load('cookie') || ''

    const data = {
      data: json,
      cookie: cookie,
    }

    axios
      .post(url, data, {
        headers: {
          'Content-Type': 'application/json',
        },
      })
      .then((response) => {
        const re: re_data = {
          statusCode: response.status,
          data: response.data,
        }
        callback(re)
        //console.log(response)
      })
      .catch((error) => {
        console.error('There was an error!', error)
      })
  },
}

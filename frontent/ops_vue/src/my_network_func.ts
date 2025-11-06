import axios from 'axios'
import { myfuncs } from './myfunc'

export const my_network_func = {
  post_json(path: string) {
    const head_path = '/api/v1'
    const url = head_path + path
    const cookie = myfuncs.load('cookie') || ''

    const data = {

      cookie: cookie,
    }

    axios
      .post(url, data, {
        headers: {
          'Content-Type': 'application/json',
        },
      })
      .then((response) => {
        console.log(response)
      })
      .catch((error) => {
        console.error('There was an error!', error)
      })
  },
}

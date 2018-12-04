#!/usr/bin/env python
import json
import urllib2
import argparse


__MOM_URL = 'https://mom.skmobi.com/'
__UUID = '19bac740-d024-0136-6e59-0a3cd2a2f8ba'


def build_parser():
    parser = argparse.ArgumentParser(description='Manage authy-pam mocking setup for CI tests')
    parser.add_argument('-r', dest='register', action='store_true',
                        help='Create MOM token')
    parser.add_argument('-t', dest='token', type=str,
                        help='MOM token to use')
    parser.add_argument('-p', dest='push', type=str, metavar='DEST',
                        help='Send push to DEST')
    parser.add_argument('-s', dest='sms', type=str, metavar='DEST',
                        help='Send SMS to DEST')
    parser.add_argument('-a', dest='approve', action='store_true',
                        help='Approve push')
    parser.add_argument('-d', dest='deny', action='store_true',
                        help='Deny push')
    return parser


def _get(path):
    req = urllib2.Request(
        __MOM_URL + path,
        headers={'Content-Type': 'application/json', 'User-Agent': 'not_urllib_due_to_CloudFlare'}
    )
    f = urllib2.urlopen(req)
    r = f.read()
    f.close()
    return json.loads(r)


def _post(path, data):
    req = urllib2.Request(
        __MOM_URL + path,
        data,
        headers={'Content-Type': 'application/json', 'User-Agent': 'not_urllib_due_to_CloudFlare'}
    )
    f = urllib2.urlopen(req)
    r = f.read()
    f.close()
    return json.loads(r)


def setup_path(token, **kw):
    return _post(
        'setup/%s/' % token,
        json.dumps(kw)
    )


def change_approval(token, status='pending'):
    body = {
        "approval_request": {
            "_app_name": "OwlBank",
            "_app_serial_id": 15861,
            "_authy_id": 245624,
            "_id": "59e8eaaa1e145e587b923fc1",
            "_user_email": "help@twilio.com",
            "app_id": "542b1b82f92ea10597000d6d",
            "created_at": "2017-10-19T18:10:50Z",
            "notified": False,
            "processed_at": "2017-10-19T18:11:37Z",
            "seconds_to_expire": 86400,
            "status": "%s" % status,
            "updated_at": "2017-10-19T18:11:37Z",
            "user_id": "145642564",
            "uuid": __UUID,
        },
        "success": True
    }
    return setup_path(
        token,
        path='onetouch/json/approval_requests/%s' % __UUID,
        status_code=200,
        body=json.dumps(body),
        content_type='application/json',
    )


def push(token, user_id):
    change_approval(token)
    return setup_path(
        token,
        path='onetouch/json/users/%s/approval_requests' % user_id,
        status_code=200,
        body='{"approval_request":{"uuid":"%s"},"success":true}' % __UUID,
        content_type='application/json',
    )


def sms(user_id):
    pass


def main():
    _p = build_parser()
    args = _p.parse_args()

    if args.register:
        print(_get('register')['token'])
        return

    if not args.token:
        _p.error('-t required')

    if args.push:
        print(push(args.token, args.push))
    elif args.sms:
        print(sms(args.token, args.push))
    elif args.approve:
        print(change_approval(args.token, 'approved'))
    elif args.deny:
        print(change_approval(args.token, 'denied'))


if __name__ == '__main__':
    main()

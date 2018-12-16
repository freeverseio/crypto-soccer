#!/usr/bin/env python
import os, sys
import subprocess
from xml.dom import minidom
import argparse
import sha3
from playerdb import *

# see : /System/Library/Frameworks/Python.framework/Versions/2.7/lib/python2.7/xml/dom/minidom.py

def is_odd(n):
    return n & 1

def get_decimal_hash(x):
    return int(sha3.sha3_256(x).hexdigest(),16)

def rgb2hex(r,g,b):
    hex = "#{:02x}{:02x}{:02x}".format(r,g,b)
    return hex

def get_svg_header(x=0, y=0, w=700, h=1000):
    return '<svg xmlns="http://www.w3.org/2000/svg" viewBox="{x} {y} {w} {h}" enable-background="new {x} {y} {w} {h}" xmlns:xlink="http://www.w3.org/1999/xlink">'.format(x=x,y=y,w=w,h=h)

def get_svg_footer():
    return '</svg>'

class SvgPath:
    fill=None
    d=None
    stroke=None
    stroke_width=None
    stroke_linecap=None
    stroke_linejoin=None

    def __init__(self, fill, d, stroke=None, stroke_width=None, stroke_linecap=None, stroke_linejoin=None):
        self.fill = fill
        self.d = d
        self.stroke=stroke
        self.stroke_width=stroke_width
        self.stroke_linecap=stroke_linecap
        self.stroke_linejoin=stroke_linejoin

    def toxml(self):
        result='\n<path fill="{p.fill}" d="{p.d}"'.format(p=self)

        if self.stroke:
            result += ' stroke="{p.stroke}"'.format(p=self)
        if self.stroke_width:
            result += ' stroke-width="{p.stroke_width}"'.format(p=self)
        if self.stroke_linecap:
            result += ' stroke-linecap="{p.stroke_linecap}"'.format(p=self)
        if self.stroke_linejoin:
            result += ' stroke-join="{p.stroke_linejoin}"'.format(p=self)

        result += '/>'

        return result

class SvgGroup:
    name=None
    transform=None
    paths=[]
    groups=[]
    display=None

    def __init__(self):
        pass

    def __init__(self, name=None, transform=None, paths=[], groups=[], display=None):
        self.name = name
        self.transform = transform
        self.display = display

        if paths:
            if not isinstance(paths, SvgPath) and not isinstance(paths, list):
                print 'Error: paths must be either an SvgPath or a list of SvgPath'
                sys.exit(-1)
            if isinstance(paths,SvgPath):
                self.paths = [paths]
            else:
                self.paths = paths
        else:
            self.paths = []

        if groups:
            if not isinstance(groups, SvgGroup) and not isinstance(groups, list):
                print 'Error: groups must be either an SvgGroup or a list of SvgGroup'
                sys.exit(-1)
            if isinstance(groups,SvgGroup):
                self.groups = [groups]
            else:
                self.groups = groups
        else:
            self.groups = []


    def toxml(self):
        result='\n<g'

        if self.name:
            result += ' id="{p.name}"'.format(p=self)
        if self.transform:
            result += ' transform="{p.transform}"'.format(p=self)
        if self.display:
            result += ' display="{p.display}"'.format(p=self)
        result+='>'

        for path in self.paths:
            result += path.toxml()
        for group in self.groups:
            result += group.toxml()

        result+='\n</g>'
        return result

def create_svg(groups, filename):
    f = open(filename + '.svg', 'w')
    f.write(get_svg_header())
    f.write('\n<g>')

    for group in groups:
        f.write(group.toxml())

    f.write('\n</g>')
    f.write(get_svg_footer())
    f.close()

def get_arms_path(n, color):
    return SvgPath(
        fill=color,
        d=arms_db[n],
        )

def get_head_path(n, color):
    return SvgPath(
        fill=color,
        d=head_db[n],
        )
def get_neck_path(n, color):
    return SvgPath(
        fill=color,
        d=neck_db[n],
        )
def get_nose_path(n, color):
    return SvgPath(fill='none',
            d=nose_db[n],
            stroke=color,
            stroke_linecap="round",
            stroke_linejoin="round",
            stroke_width="18")

def get_eyebrows_path(n, color):
    return SvgPath(fill='none',
            d=eyebrows_db[n],
            stroke=color,
            stroke_linecap="round",
            stroke_linejoin="round",
            stroke_width="9")

def get_body(n, color):
    return SvgGroup(
        name = 'body',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            get_head_path(n,color),
            get_neck_path(n, color),
            get_arms_path(n, color)
            ]
        )

def get_hair(n, color):
    return SvgGroup(
        name = 'hair',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [SvgPath(fill=color, d=hair_db[n])]
        )

def get_lips(n):
    return SvgGroup(
        name = 'lips',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            SvgPath(fill='#BD6A5B', d=lip_inf_db[n]),
            SvgPath(fill='#BD6A5B', d=lip_sup_db[n]),
            SvgPath(fill='#B05F4F', d=mouth_db[n]),
            ]
        )

def get_nose(n, color):
    return SvgGroup(
        name = 'nose',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [get_nose_path(n, color)]
        )

def get_ears(n, color):
    return SvgGroup(
        name = 'ears',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [SvgPath(fill=color, d=ears_db[n])]
        )

def get_eyebrows(n, color):
    return SvgGroup(
        name = 'eyebrows',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [get_eyebrows_path(n, color)]
        )

def get_pupils(n, color):
    return SvgGroup(
        name = 'pupils',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            SvgPath(fill=color, d=pupil_left_db[n]),
            SvgPath(fill=color, d=pupil_right_db[n]),
            ]
        )

def get_iris(n, color):
    return SvgGroup(
        name = 'iris',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            SvgPath(fill=color, d=iris_left_db[n]),
            SvgPath(fill=color, d=iris_right_db[n]),
            ]
        )

def get_teeth(n, extras_number):

    paths = [
        SvgPath(fill='#E8E8E8', d=teeth_inf_db[n]),
        SvgPath(fill='#FFFFFF', d=teeth_sup_db[n]),
        ]

    if extras_number != None:
        paths += [SvgPath(fill='#FFFFFF', d=teeth_extras_db[extras_number])]

    return SvgGroup(
        name = 'teeth',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = paths,
        )

def get_iris(n, color):
    return SvgGroup(
        name = 'iris',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            SvgPath(fill=color, d=iris_left_db[n]),
            SvgPath(fill=color, d=iris_right_db[n]),
            ]
        )

def get_tshirt(n, color):
    return SvgGroup(
        name = 'tshirt',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            SvgPath(fill=color, d=torso_db[n]),
            SvgPath(fill=color, d=sleeves_db[n]),
            SvgPath(fill=color, d=shoulders_db[n]),
            ]
        )

def get_body_color(hash_str):
    return '#' + hash_str[0:6]

def get_hair_color(hash_str):
    return '#' + hash_str[6:12]

def get_tshirt_color(hash_str):
    return '#' + hash_str[12:18]

def get_teeth_extra_type(n):
    if is_odd(n):
        return None
    return 0

def generate_player(name):
    hash_str = sha3.sha3_256(name).hexdigest()
    body_type = 0
    body_color = get_body_color(hash_str)
    hair_type = 0
    hair_color = get_hair_color(hash_str)
    lips_type = 0
    nose_color = '#B05F4F'
    nose_type = 0
    ears_type = 0
    eyebrows_type = 0
    eyebrows_color='#2A2111'
    pupils_type = 0
    pupils_color = body_color #'#66CC66'
    iris_type = 0
    iris_color = 'black'
    teeth_type = 0
    teeth_extra_type = get_teeth_extra_type(int(hash_str[0], 16))
    tshirt_type = 0
    tshirt_color = get_tshirt_color(hash_str)

    return [
            get_body(body_type, body_color),
            get_hair(hair_type, hair_color),
            get_lips(lips_type),
            get_nose(nose_type, nose_color),
            get_ears(ears_type, body_color),
            get_iris(iris_type, iris_color),
            get_pupils(pupils_type, pupils_color),
            get_eyebrows(eyebrows_type, eyebrows_color),
            get_teeth(teeth_type, teeth_extra_type),
            get_tshirt(tshirt_type, tshirt_color),
           ]


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='script to extract elements from svg or merge svgs',
            formatter_class=argparse.ArgumentDefaultsHelpFormatter,
            usage='./generate_player -n <name> -o <output_name>'
            )
    parser.add_argument('-n', '--name', help='name from which to generate a player', required=True)
    parser.add_argument('-o', '--output', help='name from which to generate a player', required=True)

    args = parser.parse_args()

    player_name = args.name
    output_name = args.output

    create_svg(generate_player(player_name), output_name)

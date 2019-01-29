#!/usr/bin/env python
import os, sys
import subprocess
#from xml.dom import minidom
import argparse
import sha3
from playerdb import *

# see : /System/Library/Frameworks/Python.framework/Versions/2.7/lib/python2.7/xml/dom/minidom.py

def is_odd(n):
    return n & 1

def clamp(val, minimum=0, maximum=255):
    if val < minimum:
        return minimum
    if val > maximum:
        return maximum
    return val

def colorscale(hexstr, scalefactor):
    """
    Scales a hex string by ``scalefactor``. Returns scaled hex string.

    To darken the color, use a float value between 0 and 1.
    To brighten the color, use a float value greater than 1.

    >>> colorscale("#DF3C3C", .5)
    #6F1E1E
    >>> colorscale("#52D24F", 1.6)
    #83FF7E
    >>> colorscale("#4F75D2", 1)
    #4F75D2
    """

    hexstr = hexstr.strip('#')

    if scalefactor < 0 or len(hexstr) != 6:
        return hexstr

    r, g, b = int(hexstr[:2], 16), int(hexstr[2:4], 16), int(hexstr[4:], 16)

    r = clamp(r * scalefactor)
    g = clamp(g * scalefactor)
    b = clamp(b * scalefactor)

    return "#%02x%02x%02x" % (r, g, b)

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
    fill_opacity=None
    d=None
    stroke=None
    stroke_width=None
    stroke_linecap=None
    stroke_linejoin=None

    def __init__(self, fill, d, fill_opacity=None, stroke=None, stroke_width=None, stroke_linecap=None, stroke_linejoin=None):
        self.fill = fill
        self.d = d
        self.fill_opacity=fill_opacity
        self.stroke=stroke
        self.stroke_width=stroke_width
        self.stroke_linecap=stroke_linecap
        self.stroke_linejoin=stroke_linejoin

    def toxml(self):
        result='\n<path fill="{p.fill}" d="{p.d}"'.format(p=self)

        if self.fill_opacity:
            result += ' fill-opacity="{p.fill_opacity}"'.format(p=self)
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

def get_shield_border_path(n, color):
    return SvgPath(fill='none',
            d=shield_border_db[n],
            stroke=color,
            stroke_linecap="round",
            stroke_linejoin="round",
            stroke_width="10.7")

def get_shield_top_path(n, color):
    return SvgPath(fill='none',
            d=shield_border_db[n],
            stroke=color,
            stroke_linecap="round",
            stroke_linejoin="round",
            stroke_width="12.75")

def get_arms(n, color):
    return SvgGroup(
        name = 'arms',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            get_arms_path(n, color)
            ]
        )

def get_shield(n, color, border_color):
    return SvgGroup(
        name = 'shield',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            SvgPath(fill=color, d=shield_filling_db[n]),
            get_shield_border_path(n, border_color),
            get_shield_top_path(n, border_color)
            ]
        )

def get_head(n, color):
    return SvgGroup(
        name = 'head',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            get_neck_path(n, color),
            SvgPath(fill='#55321B', fill_opacity=0.098, d=neck_shadow_db[n]),
            SvgPath(fill='#2C2411', fill_opacity=0.298, d=neck_side_shadow_db[n]),
            get_head_path(n,color),
            ]
        )

def get_hair(n, color):
    return SvgGroup(
        name = 'hair',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            SvgPath(fill=color, d=hair_db[n])
            ]
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

def get_tshirt_border(n, color):
    return SvgGroup(
        name = 'tshirtborder',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            SvgPath(fill=color, d=tshirt_border_db[n]),
            ]
        )

def get_shorts(n, color):
    return SvgGroup(
        name = 'shorts',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            SvgPath(fill=color, d=shorts_start_db[n]),
            SvgPath(fill=color, d=shorts_end_db[n]),
            SvgPath(fill='#2C2411', fill_opacity=0.298, d=shorts_shadow_db[n]),
            SvgPath(fill='#2C2411', fill_opacity=0.298, d=left_leg_shadow_db[n]),
            ]
        )

def get_torso_shadow(n):
    return SvgGroup(
        name = 'torso_shadow',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            SvgPath(fill='#2C2411', fill_opacity=0.298, d=torso_shadow_db[n]),
            SvgPath(fill='#2C2411', fill_opacity=0.298, d=left_arm_shadow_db[n]),
            SvgPath(fill='#2C2411', fill_opacity=0.298, d=right_arm_shadow_db[n]),
            ]
        )

def get_facial_shadow(n):
    return SvgGroup(
        name = 'facial_shadow',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [
            SvgPath(fill='#55321B', fill_opacity=0.098, d=ears_shadow_db[n]),
            SvgPath(fill='#4E3B26', fill_opacity=0.498, d=eyes_shadow_db[n]),
            SvgPath(fill='#4E3923', fill_opacity=0.498, d=mouth_shadow_db[n]),
            SvgPath(fill='#2C2411', fill_opacity=0.2, d=teeth_shadow_db[n]),
            SvgPath(fill='#2C2411', fill_opacity=0.2, d=cheeks_shadow_db[n]),
            SvgPath(fill='#2C2411', fill_opacity=0.298, d=face_shadow_db[n]),
    ]
    )

def get_body_color(hash_str):
    return '#' + hash_str[0:6]

def get_hair_color(hash_str):
    return '#' + hash_str[6:12]

def get_tshirt_color(hash_str):
    return '#' + hash_str[12:18]

def get_tshirt_border_color(hash_str):
    return '#' + hash_str[18:24]

def get_shorts_color(hash_str):
    return '#' + hash_str[24:30]

def get_iris_color(hash_str):
    return '#' + hash_str[30:36]

def get_shield_color(hash_str):
    return '#' + hash_str[36:42]

def get_shield_border_color(hash_str):
    return '#' + hash_str[42:48]

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
    eyebrows_color=colorscale(hair_color, 0.3) #'#2A2111'
    pupils_type = 0
    pupils_color = body_color #'#66CC66'
    iris_type = 0
    iris_color = get_iris_color(hash_str) #'black'
    teeth_type = 0
    teeth_extra_type = get_teeth_extra_type(int(hash_str[0], 16))
    tshirt_type = 0
    tshirt_color = get_tshirt_color(hash_str)
    tshirt_border_type = 0
    tshirt_border_color = get_tshirt_border_color(hash_str)
    shorts_type=0
    shorts_color=get_shorts_color(hash_str)
    shield_type=0
    shield_color=get_shield_color(hash_str)
    shield_border_color=get_shield_border_color(hash_str)

    return [
            get_arms(body_type, body_color),
            get_shorts(shorts_type, shorts_color),
            get_tshirt_border(tshirt_border_type, tshirt_border_color),
            get_tshirt(tshirt_type, tshirt_color),
            get_shield(shield_type, shield_color, shield_border_color),
            get_torso_shadow(body_type),
            get_head(body_type, body_color),
            get_hair(hair_type, hair_color),
            get_lips(lips_type),
            get_nose(nose_type, nose_color),
            get_ears(ears_type, body_color),
            get_iris(iris_type, iris_color),
            get_pupils(pupils_type, pupils_color),
            get_eyebrows(eyebrows_type, eyebrows_color),
            get_teeth(teeth_type, teeth_extra_type),
            get_facial_shadow(body_type),
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

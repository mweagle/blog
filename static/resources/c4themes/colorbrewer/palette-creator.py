import sys
import os
import json
import re
from string import Template
from datetime import datetime, timezone
from matplotlib.colors import ListedColormap
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt

CURRENT_UTC_TIME = datetime.now(timezone.utc)
THEME_OUTPUT_COMMON = f"""
' Metadata
' Created: {CURRENT_UTC_TIME}
"""

THEME_OUTPUT_TEMPLATE_09 = (
    """
!$$COLOR_INDEX_0="${color0}"
!$$COLOR_INDEX_1="${color1}"
!$$COLOR_INDEX_2="${color2}"
!$$COLOR_INDEX_3="${color3}"
!$$COLOR_INDEX_4="${color4}"
!$$COLOR_INDEX_5="${color5}"
!$$COLOR_INDEX_6="${color6}"
!$$COLOR_INDEX_7="${color7}"
!$$COLOR_INDEX_8="${color8}"

!$$COLOR_BORDER_COLOR="#000000"

!$$COLOR_TEXT_LIGHT="${color1}"
!$$COLOR_TEXT_NEUTRAL="${color4}"
!$$COLOR_TEXT_DARK="${color8}"

!$$COLOR_REL_COLOR="${color8}"
"""
    + THEME_OUTPUT_COMMON
)

THEME_OUTPUT_TEMPLATE_11 = (
    """
!$$COLOR_INDEX_0="${color0}"
!$$COLOR_INDEX_1="${color1}"
!$$COLOR_INDEX_2="${color2}"
!$$COLOR_INDEX_3="${color3}"
!$$COLOR_INDEX_4="${color4}"
!$$COLOR_INDEX_5="${color5}"
!$$COLOR_INDEX_6="${color6}"
!$$COLOR_INDEX_7="${color7}"
!$$COLOR_INDEX_8="${color8}"
!$$COLOR_INDEX_9="${color9}"
!$$COLOR_INDEX_10="${color10}"

!$$COLOR_BORDER_COLOR="#000000"

!$$COLOR_TEXT_LIGHT="${color1}"
!$$COLOR_TEXT_NEUTRAL="${color5}"
!$$COLOR_TEXT_DARK="${color9}"

!$$COLOR_REL_COLOR="${color10}"
"""
    + THEME_OUTPUT_COMMON
)


def rgb_str_to_hex(r, g, b):
    return ("{:02X}" * 3).format(int(r), int(g), int(b))


# Read the input file
theme_template_09 = Template(THEME_OUTPUT_TEMPLATE_09)
theme_template_11 = Template(THEME_OUTPUT_TEMPLATE_11)

with open("./colorbrewer.json", "r") as file:
    jsonData = json.load(file)
    for each_key, each_map in jsonData.items():
        palette_type = each_map["type"]
        map_colors = {}
        for each_list_count_key, each_list_colors in each_map.items():
            if each_list_count_key != "9" and each_list_count_key != "11":
                continue

            hex_colors = []
            for each_color in each_list_colors:
                # Extract the string
                re_result = re.search("rgb\((\d+)\w*,(\d+)\w*,(\d+)\w*\)", each_color)
                hex_color = rgb_str_to_hex(
                    re_result.group(1), re_result.group(2), re_result.group(3)
                )
                hex_colors.append("#{0}".format(hex_color))

            output_data = ""
            if len(hex_colors) == 9:
                # Great - take the name, create a file with the template content, call it good
                output_data = theme_template_09.substitute(
                    color0=hex_colors[0],
                    color1=hex_colors[1],
                    color2=hex_colors[2],
                    color3=hex_colors[3],
                    color4=hex_colors[4],
                    color5=hex_colors[5],
                    color6=hex_colors[6],
                    color7=hex_colors[7],
                    color8=hex_colors[8],
                )
            elif len(hex_colors) == 11:
                # Great - take the name, create a file with the template content, call it good
                output_data = theme_template_11.substitute(
                    color0=hex_colors[0],
                    color1=hex_colors[1],
                    color2=hex_colors[2],
                    color3=hex_colors[3],
                    color4=hex_colors[4],
                    color5=hex_colors[5],
                    color6=hex_colors[6],
                    color7=hex_colors[7],
                    color8=hex_colors[8],
                    color9=hex_colors[9],
                    color10=hex_colors[10],
                )
            if output_data == "":
                raise RuntimeError("Invalid output template created")

            # Write-out the data
            output_basename = "./palettes/{0}-{1}-{2}".format(
                palette_type, each_key, len(hex_colors)
            )
            output_puml_filename = "{0}.puml".format(output_basename)
            output_png_filename = "{0}.png".format(output_basename)
            with open(output_puml_filename, "w") as theme_file_output:
                theme_file_output.write(output_data)
            print("Created theme: {0}".format(output_puml_filename))
            # Create the thumbnail...
            seaborn_pallete = sns.color_palette(hex_colors)
            sns.palplot(seaborn_pallete)
            plt.savefig(output_png_filename)
            plt.close("all")

            print("Created thumbnail: {0}".format(output_png_filename))

<?php

class ScmStatus
{
	private static $parsed = false;
	private static $filepath = "REVISION.json";

	public static function setFilepath($filepath = "REVISION.json")
	{
		self::$filepath = realpath($filepath);
	}

	private static function parse()
	{
		if (self::$parsed)
		{
			return self::$parsed;
		}

		$info = @file_get_contents(self::$filepath);

		if (!$info)
		{
			return false;
		}

		self::$parsed = json_decode($info, true);

		return self::$parsed;
	}

	public static function format($format_str, array $options = array())
	{
		$info = self::parse();

		if (isset($options['return_if_fail']))
		{
			$return_if_fail = $options['return_if_fail'];
		}
		else
		{
			$return_if_fail = "";
		}

		if (!$info)
		{
			return $return_if_fail;
		}

		$tokens = array(
			"%T" => 'type',
			"%n" => 'dec',
			"%r" => 'hex.short',
			"%R" => 'hex.full',
			"%b" => 'branch',
			"%t" => 'tags',
			"%a" => 'author.name',
			"%e" => 'author.email',
			"%m" => 'message',
			"%s" => 'subject',
		);

		if (isset($info['commit_timestamp']))
		{
			$tokens["%d"] = 'commit_date';
			$tokens["%U"] = 'commit_timestamp';
		}
		else
		{
			$tokens["%d"] = 'commit_date.date';
			$tokens["%U"] = 'commit_date.timestamp';
		}

		if (isset($options['format_date']) && isset($info['commit_timestamp']))
		{
			$info['commit_date_formatted'] = date($options['format_date'], $info['commit_timestamp']);
			$tokens['%F'] = 'commit_date_formatted';
		}
		elseif (isset($options['format_date']))
		{
			$info['commit_date_formatted'] = date($options['format_date'], $info['commit_date']['timestamp']);
			$tokens['%F'] = 'commit_date_formatted';
		}

		if (isset($options['delimiter']))
		{
			$delimiter = $options['delimiter'];
		}
		else
		{
			$delimiter = ",";
		}

		$tbr = $format_str;

		foreach ($tokens as $token => $key)
		{
			$key = explode(".", $key);

			$val = $info;
			for ($i = 0; $i < count($key); $i++)
			{
				if (isset($val[$key[$i]]))
				{
					$val = $val[$key[$i]];
				}
				else
				{
					$val = false;
				}
			}

			if (is_array($val))
			{
				$val = implode($delimiter, $val);
			}

			if (!empty($val))
			{
				$tbr = str_replace($token, $val, $tbr);
			}
		}

		return $tbr;
	}

	public static function load()
	{
		return self::parse();
	}
}
